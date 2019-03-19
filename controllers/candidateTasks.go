package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)

func GetCandidateInfo(c *gin.Context) {
    roll := c.Param("roll")
    
    candidate, err := ElectionDb.GetCandidate(roll)
    if err != nil {
        c.String(http.StatusNotFound, "This student is not contesting elections.")
        return
    }
    
    simplifiedCandidate := candidate.Simplify()
    
    c.JSON(http.StatusOK, simplifiedCandidate)
}

func GetCandidateCard(c *gin.Context) {
    roll := c.Param("roll")
    
    card := "<div class=\"user_card mt-5 mb-3\" style=\"width:230px\">"
    card += "<div class=\"d-flex justify-content-center\" align=\"center\">"
    card += "<div class=\"image_container\">"
    card += "<img id=\"photo\" src=\"logo.jpg\" class=\"round_image\" alt=\"Candidate's Photo\">"
    card += "</div>"
    card += "<div class=\"d-flex justify-content-center form_container\">"
    card += "<div class=\"card-body\">"
    card += "<h3 id=\"candidateName\">%s</h3><br>"
    card += "<a target=\"_blank\" id=\"manifestoLink\" class=\"btn btn-light\" href=\"%s\">Manifesto</a><br><br>"
    card += "<input id=\"voteButton\" type=\"button\" postid=\"%s\" pubkey=\"%s\" class=\"btn btn-primary\" value=\"1st Preference\"/>"
    card += "</div></div></div></div><pre>  </pre>"
        
    candidate, err := ElectionDb.GetCandidate(roll)
    if err != nil {
        c.String(http.StatusNotFound, "Candidate Not Found");
        return
    }
    
    formattedCard := fmt.Sprintf(card,candidate.Name,candidate.Manifesto, candidate.PostID, candidate.PublicKey)
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedCard))
}
