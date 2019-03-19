package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
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
