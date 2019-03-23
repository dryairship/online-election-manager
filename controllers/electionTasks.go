package controllers

import (
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

func SubmitVote(c *gin.Context) {
    roll, err := utils.GetSessionID(c)
    if err != nil {
        c.String(http.StatusForbidden, "Unauthorized vote submission.")
        return
    }
    
    voter, err := ElectionDb.FindVoter(roll)
    if err != nil {
        c.String(http.StatusNotFound, "Voter not found.")
        return
    }
    
    if voter.Voted {
        c.String(http.StatusForbidden, "This voter has already voted.")
        return
    }
    
    var receivedVotes []models.ReceivedVote
    err = c.BindJSON(&receivedVotes)
    if err != nil {
        c.String(http.StatusBadRequest, "Data format not recognized.")
        return
    }
    
    ballotID := make([]models.BallotID, len(receivedVotes))
    for i, receivedVote := range receivedVotes {
        ballotID[i] = receivedVote.GetBallotID()
        vote := receivedVote.GetVote()
        err = ElectionDb.InsertVote(&vote)
        if err != nil {
            c.String(http.StatusInternalServerError, "Database Error")
            return
        }
    }
    
    newVoter := models.Voter{
        Roll:       voter.Roll,
        Name:       voter.Name,
        Email:      voter.Email,
        Password:   voter.Password,
        AuthCode:   voter.AuthCode,
        BallotID:   ballotID,
        Voted:      true,
    }
    
    err = ElectionDb.UpdateVoter(roll, &newVoter)
    if err != nil {
        c.String(http.StatusInternalServerError, "Database Error")
        return
    }
    
    c.String(http.StatusOK, "Votes successfully submitted.")
}
