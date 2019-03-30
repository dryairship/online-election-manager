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

func RegisterCEO(c *gin.Context) {
    passHash := c.PostForm("pass")
    authCode := c.PostForm("auth")
    
    ceo, err := ElectionDb.GetCEO()
    if err != nil {
        c.String(http.StatusForbidden, "You need to get a verification<br>mail before you register.")
        return
    }
    
    if ceo.AuthCode == "" {
        c.String(http.StatusForbidden, "CEO has already registered.")
        return
    }
    
    if ceo.AuthCode != authCode {
        c.String(http.StatusBadRequest, "Wrong authentication code.")
        return
    }
    
    ceo.Password = passHash
    ceo.AuthCode = ""
    
    err = ElectionDb.UpdateCEO(&ceo)
    if err != nil {
        c.String(http.StatusInternalServerError, "Database Error")
    }else{
        c.String(http.StatusAccepted, "CEO successfully registered.")
    }
}

func CEOLogin(c *gin.Context) {
    passHash := c.PostForm("pass")
    ceo, err := ElectionDb.GetCEO()
    if err != nil {
        c.String(http.StatusForbidden, "CEO has not yet registered.")
        return
    }
    
    if ceo.Password != passHash {
        c.String(http.StatusForbidden, "Invalid Password.")
        return
    }
    
    utils.StartSession(c)
    
    c.JSON(http.StatusOK, &ceo)
}

