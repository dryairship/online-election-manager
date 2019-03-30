package controllers

import (
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/dryairship/online-election-manager/db"
    "github.com/gin-gonic/gin"
    "net/http"
    "regexp"
)

var ElectionDb db.ElectionDatabase
var Posts      []models.Post

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
            return false, &UserError{"Verification mail has already<br>been sent to this student."}
        }
    }
}

func GetPostsForVoter(roll string) []models.VotablePost {
    votablePosts := []models.VotablePost{}
    for _, post := range Posts {
        pattern := post.VoterRegex
        canVote, err := regexp.MatchString(pattern, roll)
        if err == nil && canVote {
            votablePosts = append(votablePosts, post.ConvertToVotablePost())
        }
    }
    return votablePosts
}

func GetVotablePosts(c *gin.Context){
    roll := c.Param("roll")
    
    id, err := utils.GetSessionID(c)
    if err != nil || id != roll {
        c.String(http.StatusForbidden, "You are not authorized to use this API.")
        return
    }
    
    votablePosts := GetPostsForVoter(roll)
    c.JSON(http.StatusOK, votablePosts)
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
    
    if passHash == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
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
    
    if roll == "CEO" {
        SendMailToCEO(c)
        return
    }
    
    _, err := CanMailBeSentToStudent(roll)
    if err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    
    skeleton, err := ElectionDb.FindStudentSkeleton(roll)
    if err != nil {
        c.String(http.StatusNotFound, "Invalid Roll Number.<br>Student does not exist.")
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
    c.String(http.StatusAccepted, "Verification Mail successfully sent<br>to "+voter.Email)
}

func CheckUserLogin(c *gin.Context) {
    roll := c.PostForm("roll")
    passHash := c.PostForm("pass")
    
    voter, err := ElectionDb.FindVoter(roll)
    if err != nil {
        c.String(http.StatusForbidden, "This student has not registered.")
        return
    }
    
    if voter.Password != passHash {
        c.String(http.StatusForbidden, "Invalid Password.")
        return
    }
    
    utils.StartSession(c)
    
    simplifiedVoter := voter.Simplify()
    c.JSON(http.StatusOK, &simplifiedVoter)
}

