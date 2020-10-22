package controllers

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

var ElectionDb db.ElectionDatabase
var Posts []models.Post

type UserError struct {
	reason string
}

type CAPTCHA struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

func (err *UserError) Error() string {
	return err.reason
}

// Function to check if a verification mail can be sent to the student.
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

// Function to match voter's roll number with the regular expression for a post
// to determine whether or not the voter is eligible to vote for a post.
func GetPostsForVoter(roll string) []models.VotablePost {
	votablePosts := []models.VotablePost{}
	for _, post := range Posts {
		pattern := post.VoterRegex
		canVote, err := regexp.MatchString(pattern, roll)
		if err == nil && canVote && !post.Resolved {
			votablePosts = append(votablePosts, post.ConvertToVotablePost())
		}
	}
	return votablePosts
}

// API handler to fetch all posts for which this voter is eligible to vote.
func GetVotablePosts(c *gin.Context) {
	roll, err := utils.GetSessionID(c)
	if err != nil {
		log.Println("[WARN] Unauthorized getVotablePosts attempt: ", err.Error())
		c.String(http.StatusForbidden, "You are not authorized to use this API.")
		return
	}

	votablePosts := GetPostsForVoter(roll)
	c.JSON(http.StatusOK, votablePosts)
}

func GetAllPosts(c *gin.Context) {
	c.JSON(http.StatusOK, Posts)
}

// API handler to register a new voter.
func RegisterNewVoter(c *gin.Context) {
	if config.ElectionState == config.VotingStopped {
		log.Println("[WARN] Registration attempt after end of registration period: ", c.PostForm("roll"))
		c.String(http.StatusForbidden, "Registration period is over.")
		return
	}

	roll := c.PostForm("roll")
	passHash := c.PostForm("pass")
	authCode := c.PostForm("auth")

	if roll == "CEO" {
		RegisterCEO(c)
		return
	}

	if roll[0] == 'P' {
		RegisterCandidate(c)
		return
	}

	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		c.String(http.StatusForbidden, "You need to get a verification mail before you register.")
		return
	}

	if voter.AuthCode == "" {
		log.Println("[WARN] Re-registration attempt: ", voter)
		c.String(http.StatusForbidden, "Student has already registered.")
		return
	}

	if voter.AuthCode != authCode {
		c.String(http.StatusBadRequest, "Wrong authentication code.")
		return
	}

	if passHash == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		log.Println("[WARN] Registration attempt with empty password: ", voter)
		c.String(http.StatusBadRequest, "Empty password is not allowed.")
		return
	}

	voter.Password = passHash
	voter.AuthCode = ""

	err = ElectionDb.UpdateVoter(roll, &voter)
	if err != nil {
		log.Println("[ERROR] Database error while registering voter: ", voter)
		c.String(http.StatusInternalServerError, "Database Error")
	} else {
		c.String(http.StatusAccepted, "Voter successfully registered.")
	}
}

// API handler to send a verification mail to the student.s
func SendMailToStudent(c *gin.Context) {
	if config.ElectionState == config.VotingStopped {
		log.Println("[WARN] Verification-mail attempt after end of registration period: ", c.Param("roll"))
		c.String(http.StatusForbidden, "Registration period is over.")
		return
	}

	var captcha CAPTCHA
	err := c.BindJSON(&captcha)
	if err != nil {
		log.Println("[ERROR] CAPTCHA data failed to bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Invalid data format")
		return
	}

	success := utils.VerifyCaptcha(captcha.Id, captcha.Value)
	if !success {
		c.String(http.StatusBadRequest, "Incorrect CAPTCHA")
		return
	}

	roll := c.Param("roll")

	if roll == "CEO" {
		SendMailToCEO(c)
		return
	}

	if roll[0] == 'P' {
		SendMailToCandidate(c)
		return
	}

	_, err = CanMailBeSentToStudent(roll)
	if err != nil {
		log.Println("[WARN] Verification-mail not sent to voter: ", roll, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	skeleton, err := ElectionDb.FindStudentSkeleton(roll)
	if err != nil {
		log.Println("[ERROR] Student not found in database while sending verification mail: ", roll, err.Error())
		c.String(http.StatusNotFound, "Invalid Roll Number. Student does not exist.")
		return
	}

	voter := skeleton.CreateVoter(utils.GetRandomAuthCode())
	recipient := voter.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "a voter")
	if err != nil {
		log.Println("[ERROR] Mailer not working while sending mail to voter: ", voter, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	} else {
		err = ElectionDb.AddNewVoter(&voter)
		if err != nil {
			log.Println("[ERROR] Database error when sending mail to voter: ", voter, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}
	}
	c.String(http.StatusAccepted, "Verification Mail successfully sent to "+voter.Email)
}

// API handler to check the user's login credentials.
func CheckUserLogin(c *gin.Context) {
	roll := c.PostForm("roll")
	passHash := c.PostForm("pass")

	if roll == "CEO" {
		CEOLogin(c)
		return
	}

	if roll[0] == 'P' {
		CandidateLogin(c)
		return
	}

	if config.ElectionState == config.VotingNotYetStarted {
		c.String(http.StatusForbidden, "Voting has not yet started.")
		return
	}

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

func GetCaptcha(c *gin.Context) {
	id, value := utils.CreateCaptcha()
	captcha := CAPTCHA{
		Id:    id,
		Value: value,
	}

	c.JSON(http.StatusOK, &captcha)
}
