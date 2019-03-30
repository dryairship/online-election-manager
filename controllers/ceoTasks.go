package controllers

import (
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

func SendMailToCEO(c *gin.Context) {
    ceo, err := ElectionDb.GetCEO()
    if err == nil {
        c.String(http.StatusForbidden, "Verification mail has already<br>been sent to CEO.")
        return
    }
    
    skeleton, err := ElectionDb.FindStudentSkeleton(config.RollNumberOfCEO)
    if err != nil {
        c.String(http.StatusInternalServerError, "No CEO assigned.")
        return
    }
    
    ceo = skeleton.CreateCEO(utils.GetRandomAuthCode())
    recipient := ceo.GetMailRecipient()
    err = utils.SendMailTo(&recipient)
    if err != nil {
        c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
        return
    } else {
        err = ElectionDb.InsertCEO(&ceo)
        if err != nil {
            c.String(http.StatusInternalServerError, "Database error.")
            return
        }
    }
    c.String(http.StatusAccepted, "Verification Mail successfully sent<br>to "+ceo.Email)
}
