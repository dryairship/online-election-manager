package controllers

import (
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)

func GetCandidateInfo(c *gin.Context) {
    _, err := utils.GetSessionID(c)
    if err != nil {
        c.String(http.StatusForbidden, "Unauthorized request.")
        return
    }
    
    username := c.Param("username")
    
    candidate, err := ElectionDb.GetCandidate(username)
    if err != nil {
        c.String(http.StatusNotFound, "This student is not contesting elections.")
        return
    }
    
    simplifiedCandidate := candidate.Simplify()
    
    c.JSON(http.StatusOK, simplifiedCandidate)
}

func GetCandidateCard(c *gin.Context) {
    _, err := utils.GetSessionID(c)
    if err != nil {
        c.String(http.StatusForbidden, "Unauthorized request.")
        return
    }
    
    username := c.Param("username")
    
    card := "<div class=\"user_card my-5 mx-3\" style=\"width:230px\">"
    card += "<div class=\"d-flex justify-content-center\" align=\"center\">"
    card += "<div class=\"image_container\">"
    card += "<img id=\"photo\" src=\"https://oa.cc.iitk.ac.in/Oa/Jsp/Photo/%s_0.jpg\" class=\"round_image\" alt=\"Candidate's Photo\">"
    card += "</div>"
    card += "<div class=\"d-flex justify-content-center form_container\">"
    card += "<div class=\"card-body\">"
    card += "<h3 id=\"candidateName\">%s</h3><br>"
    card += "<a target=\"_blank\" id=\"manifestoLink\" class=\"btn btn-light\" href=\"%s\">Manifesto</a><br><br>"
    card += "<input id=\"voteButton\" type=\"button\" postid=\"%s\" pubkey=\"%s\" class=\"btn btn-primary\" value=\"1st Preference\"/>"
    card += "</div></div></div></div>"
        
    candidate, err := ElectionDb.GetCandidate(username)
    if err != nil {
        c.String(http.StatusNotFound, "Candidate Not Found");
        return
    }
    
    formattedCard := fmt.Sprintf(card, candidate.Roll, candidate.Name, candidate.Manifesto, candidate.PostID, candidate.PublicKey)
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedCard))
}


func SendMailToCandidate(c *gin.Context) {
    if config.ElectionState != config.VotingNotYetStarted {
        c.String(http.StatusForbidden, "Registration period for<br>candidates is over.")
        return
    }
    
    username := c.Param("roll")
    candidate, err := ElectionDb.GetCandidate(username)
    if err != nil {
        c.String(http.StatusNotFound, "This candidate is not<br>contesting elections.")
        return
    }
    
    if candidate.Password != "" {
        c.String(http.StatusForbidden, "This candidate has already<br>registered.")
        return
    } else if candidate.AuthCode != "" {
        c.String(http.StatusForbidden, "Verification mail has already<br>been sent to this candidate.")
        return
    }
    
    candidate.AuthCode = utils.GetRandomAuthCode()
    recipient := candidate.GetMailRecipient()
    err = utils.SendMailTo(&recipient)
    if err != nil {
        c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
        return
    } else {
        err = ElectionDb.UpdateCandidate(username, &candidate)
        if err != nil {
            c.String(http.StatusInternalServerError, "Database error.")
            return
        }
    }
    c.String(http.StatusAccepted, "Verification Mail successfully sent<br>to "+candidate.Email)
}

func RegisterCandidate(c *gin.Context) {
    username := c.PostForm("roll")
    passHash := c.PostForm("pass")
    authCode := c.PostForm("auth")
    
    candidate, err := ElectionDb.GetCandidate(username)
    
    if err != nil {
        c.String(http.StatusNotFound, "Candidate not found.")
        return
    }
    
    if candidate.Password != "" {
        c.String(http.StatusForbidden, "Candidate has already registered.")
        return
    } else if candidate.AuthCode == "" {
        c.String(http.StatusForbidden, "You need to get a verification<br>mail before you register.")
        return
    }
    
    if candidate.AuthCode != authCode {
        c.String(http.StatusBadRequest, "Wrong authentication code.")
        return
    }
    
    if passHash == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
        c.String(http.StatusBadRequest, "Empty password is not allowed.")
        return
    }
    
    candidate.Password = passHash
    candidate.AuthCode = ""
    
    err = ElectionDb.UpdateCandidate(username, &candidate)
    if err != nil {
        c.String(http.StatusInternalServerError, "Database Error")
    }else{
        c.String(http.StatusAccepted, "Candidate successfully registered.")
    }
}


func CandidateLogin(c *gin.Context) {
    username := c.PostForm("roll")
    passHash := c.PostForm("pass")
    candidate, err := ElectionDb.GetCandidate(username)
    if err != nil {
        c.String(http.StatusForbidden, "This candidate has not yet registered.")
        return
    }
    
    if candidate.Password != passHash {
        c.String(http.StatusForbidden, "Invalid Password.")
        return
    }
    
    utils.StartSession(c)
    
    simplifiedCandidate := candidate.Simplify()
    c.JSON(http.StatusOK, &simplifiedCandidate)
}

