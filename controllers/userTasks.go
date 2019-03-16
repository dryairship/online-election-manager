package controllers

import (
    "github.com/dryairship/online-election-manager/utils"
    "github.com/dryairship/online-election-manager/db"
    "github.com/gin-gonic/gin"
    "net/http"
)

var ElectionDb db.ElectionDatabase

type UserError struct {
    reason string
}

func (err *UserError) Error() string {
    return err.reason
}

func CanMailBeSentToStudent(roll string) (bool, error) {
    voter, err := ElectionDb.FindVoter(roll)
    if err != nil {
        return true, nil
    } else {
        if voter.AuthCode == "" {
            return false, &UserError{"Student has already registered."}
        } else {
            return false, &UserError{"Verification mail has already been sent to this student."}
        }
    }
}

func RegisterNewVoter(c *gin.Context){
    roll := c.PostForm("roll")
    passHash := c.PostForm("pass")
    authCode := c.PostForm("auth")
    
    voter, err := ElectionDb.FindVoter(roll)
    if err != nil {
        c.String(http.StatusForbidden, "You need to get a verification<br>mail before you register.")
        return
    }
    
    if voter.AuthCode == "" {
        c.String(http.StatusForbidden, "Student has already registered.")
        return
    }
    
    if voter.AuthCode != authCode {
        c.String(http.StatusBadRequest, "Wrong authentication code.")
        return
    }
    
    if passHash == "" {
        c.String(http.StatusBadRequest, "Empty password is not allowed.")
        return
    }
    
    voter.Password = passHash
    voter.AuthCode = ""
    
    err = ElectionDb.UpdateVoter(roll,&voter)
    if err != nil {
        c.String(http.StatusInternalServerError, "Database Error")
    }else{
        c.String(http.StatusAccepted, "Voter successfully registered.")
    }
}

func SendMailToStudent(c *gin.Context) {
    roll := c.Param("roll")
    _, err := CanMailBeSentToStudent(roll)
    if err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    
    skeleton, err := ElectionDb.FindStudentSkeleton(roll)
    if err != nil {
        c.String(http.StatusNotFound, "Invalid Roll Number. Student does not exist.")
        return
    }
    
    voter := skeleton.CreateVoter(utils.GetRandomAuthCode())
    recipient := voter.GetMailRecipient()
    err = utils.SendMailTo(&recipient)
    if err != nil {
        c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
        return
    } else {
        err = ElectionDb.AddNewVoter(&voter)
        if err != nil {
            c.String(http.StatusInternalServerError, "Database error.")
            return
        }
    }
    c.String(http.StatusAccepted, "Verification Mail successfully sent to "+voter.Email)
}

func SessionLogin(c *gin.Context){
    if c.PostForm("roll") == "180561" {
        c.JSON(200, gin.H{
            "success": 1,
            "roll": c.PostForm("roll"),
            "pass": c.PostForm("password"),
        })
    } else {
        c.JSON(200, gin.H{
            "success": 0,
        })
    }
}
