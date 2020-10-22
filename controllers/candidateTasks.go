package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

// Struct to accept new keys from the client.
type Keys struct {
	PublicKey  string `json:"pubkey"`
	PrivateKey string `json:"privkey"`
}

var card = "<div class=\"user_card my-5 mx-3\" style=\"width:230px\">" +
	"<div class=\"d-flex justify-content-center\" align=\"center\">" +
	"<div class=\"image_container\">" +
	"<img id=\"photo\" src=\"https://oa.cc.iitk.ac.in/Oa/Jsp/Photo/%s_0.jpg\" class=\"round_image\" alt=\"Candidate's Photo\">" +
	"</div>" +
	"<div class=\"d-flex justify-content-center form_container\">" +
	"<div class=\"card-body\">" +
	"<h3 id=\"candidateName\">%s</h3><br>" +
	"<a target=\"_blank\" id=\"manifestoLink\" class=\"btn btn-light\" href=\"%s\">Manifesto</a><br><br>" +
	"<input id=\"voteButton\" type=\"button\" postid=\"%s\" pubkey=\"%s\" class=\"btn btn-primary\" value=\"1st Preference\"/>" +
	"</div></div></div></div>"

// API handler to get information about a candidate.
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

// API handler to get the HTML card of a candidate.
// *Hacky solution* to already load the data into the HTML template.
func GetCandidateCard(c *gin.Context) {
	_, err := utils.GetSessionID(c)
	if err != nil {
		c.String(http.StatusForbidden, "Unauthorized request.")
		return
	}

	username := c.Param("username")

	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "Candidate Not Found")
		return
	}

	formattedCard := fmt.Sprintf(card, candidate.Roll, candidate.Name, candidate.Manifesto, candidate.PostID, candidate.PublicKey)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedCard))
}

// API handler to send verification mail to a candidate.
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
	err = utils.SendMailTo(&recipient, "a candidate")
	if err != nil {
		log.Println("[ERROR] Mailer error while sending mail to candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	}

	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while sending mail to candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database error.")
		return
	}

	c.String(http.StatusAccepted, "Verification Mail successfully sent<br>to "+candidate.Email)
}

// API handler to register a new candidate.
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
		log.Println("[ERROR] Database error while registering candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
	} else {
		c.String(http.StatusAccepted, "Candidate successfully registered.")
	}
}

// API handler to check candidate's login credentials.
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

// API handler to accept the unencrypted private key of a candidate.
func DeclarePrivateKey(c *gin.Context) {
	username, err := utils.GetSessionID(c)
	if err != nil {
		c.String(http.StatusForbidden, "You need to be logged in.")
		return
	}

	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "Candidate not found.")
		return
	}

	if candidate.KeyState != config.KeysGenerated {
		c.String(http.StatusBadRequest, "Candidate has already declared private key.")
		return
	}

	keys := Keys{}
	err = c.BindJSON(&keys)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	candidate.PublicKey = keys.PublicKey
	candidate.PrivateKey = keys.PrivateKey
	candidate.KeyState = config.KeysDeclared
	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while declaring private keys of candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	c.JSON(http.StatusOK, "Public Key succesfully received.")
}

// API handler to confirm candidature by updating the public and private keys of a candidate.
func ConfirmCandidature(c *gin.Context) {
	username, err := utils.GetSessionID(c)
	if err != nil {
		c.String(http.StatusForbidden, "You need to be logged in.")
		return
	}

	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "Candidate not found.")
		return
	}

	if candidate.KeyState != config.KeysNotGenerated {
		c.String(http.StatusBadRequest, "Candidate has already registered.")
		return
	}

	keys := Keys{}
	err = c.BindJSON(&keys)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	candidate.PublicKey = keys.PublicKey
	candidate.PrivateKey = keys.PrivateKey
	candidate.KeyState = config.KeysGenerated
	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while confirming candidature: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	c.JSON(http.StatusOK, "Candidature confirmed.")
}
