package controllers

import (
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
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
        go func() {
            time.Sleep(utils.GetRandomTimeDelay())
            ElectionDb.InsertVote(&vote)
        }()
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
    
    c.JSON(http.StatusOK, "Votes successfully submitted.")
}

func GetElectionState(c *gin.Context) {
    _, err := utils.GetSessionID(c)
    if err != nil {
        c.String(http.StatusForbidden, "You need to be logged in.")
        return
    }
    
    c.String(http.StatusOK, string(48+config.ElectionState))
}

