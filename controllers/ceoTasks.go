package controllers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

// Struct to accept the keys from the client.
type NewCEOData struct {
	PublicKey  string `json:"pubkey"`
	PrivateKey string `json:"privkey"`
}

// Struct to return candidates' data to the CEO.
type CandidateData struct {
	Name       string
	Roll       string
	PostID     string
	PrivateKey string
}

type CandidateResult struct {
	Roll      string   `json:"roll"`
	Name      string   `json:"name"`
	Count     int32    `json:"count"`
	BallotIDs []string `json:"ballotIds"`
	Status    string   `json:"status"`
}

type PostResult struct {
	ID         string            `json:"postid"`
	Name       string            `json:"postname"`
	Resolved   bool              `json:"resolved"`
	Candidates []CandidateResult `json:"candidates"`
}

type ResultData struct {
	Posts []PostResult `json:"posts"`
}

// API handler to send verification mail to CEO.
func SendMailToCEO(c *gin.Context) {
	ceo, err := ElectionDb.GetCEO()
	if err == nil {
		c.String(http.StatusForbidden, "Verification mail has already been sent to CEO.")
		return
	}

	skeleton, err := ElectionDb.FindStudentSkeleton(config.RollNumberOfCEO)
	if err != nil {
		c.String(http.StatusInternalServerError, "No CEO assigned.")
		return
	}

	ceo = skeleton.CreateCEO(utils.GetRandomAuthCode())
	recipient := ceo.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "the CEO")
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
	c.String(http.StatusAccepted, "Verification Mail successfully sent to "+ceo.Email)
}

// API handler to register CEO's account.
func RegisterCEO(c *gin.Context) {
	passHash := c.PostForm("pass")
	authCode := c.PostForm("auth")

	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		c.String(http.StatusForbidden, "You need to get a verification mail before you register.")
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
	} else {
		c.String(http.StatusAccepted, "CEO successfully registered.")
	}
}

// API handler to check CEO's login credentials.
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

// API handler to fetch all votes from the database.
func FetchVotes(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	votes, err := ElectionDb.GetVotes()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while fetching votes.")
		return
	}
	c.JSON(http.StatusOK, &votes)
}

func isCandidateStillInRace(postId, username string) bool {
	for _, post := range Posts {
		if post.PostID == postId && !post.Resolved {
			for _, candidate := range post.Candidates {
				if candidate == username {
					return true
				}
			}
		}
	}
	return false
}

// API handler to fetch all candidates from the database.
func FetchCandidates(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	candidates, err := ElectionDb.GetAllCandidates()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while fetching candidates.")
	}

	var data []CandidateData
	for _, candidate := range candidates {
		if isCandidateStillInRace(candidate.PostID, candidate.Username) {
			data = append(data, CandidateData{
				Name:       candidate.Name,
				Roll:       candidate.Roll,
				PostID:     candidate.PostID,
				PrivateKey: candidate.PrivateKey,
			})
		}
	}
	c.JSON(http.StatusOK, data)
}

// API handler to fetch all posts from the database.
func FetchPosts(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	posts, err := ElectionDb.GetPostsForCEO()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while fetching posts.")
		return
	}
	c.JSON(http.StatusOK, &posts)
}

// API handler to start voting process by accepting CEO's public and private keys.
func StartVoting(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	if config.ElectionState != config.VotingNotYetStarted {
		c.String(http.StatusBadRequest, "You are too late. Voting has already started.")
		return
	}

	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		c.String(http.StatusInternalServerError, "CEO not yet assigned.")
		return
	}
	newData := NewCEOData{}
	err = c.BindJSON(&newData)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	ceo.PublicKey = newData.PublicKey
	ceo.PrivateKey = newData.PrivateKey
	err = ElectionDb.UpdateCEO(&ceo)
	if err != nil {
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	config.PublicKeyOfCEO = ceo.PublicKey
	config.ElectionState = config.AcceptingVotes

	c.String(http.StatusOK, "Voting Started.")
}

// API handler to stop voting process.
func StopVoting(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	if config.ElectionState != config.AcceptingVotes {
		c.String(http.StatusBadRequest, "The system is already not accepting new votes.")
		return
	}

	config.ElectionState = config.VotingStopped
	c.String(http.StatusOK, "Voting Stopped")
}

func ResultProgress(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("%.3f%%", config.ResultProgress))
}

func GetResult(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	results, err := ElectionDb.FindAllResults()
	if err != nil {
		c.String(http.StatusBadRequest, "Cannot find results in the database.")
		return
	}

	numericResults := make([]models.NumericResult, len(results))
	for i, result := range results {
		numericResults[i] = result.ToNumericResult()
	}

	sort.Sort(models.ResultSorter(numericResults))
	c.JSON(http.StatusOK, &numericResults)
}

func PrepareForNextRound(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	err = ElectionDb.MarkAllVotersUnvoted()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	err = ElectionDb.DeleteAllVotes()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	config.ElectionState = config.VotingNotYetStarted

	c.String(http.StatusOK, "Ready for next round.")
}

func SubmitSingleVoteResults(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	var results ResultData
	err = c.BindJSON(&results)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, post := range results.Posts {
		var ballotIds []models.UsedBallotID
		postId := post.ID
		if post.Resolved {
			if err := ElectionDb.MarkPostResolved(postId); err != nil {
				c.String(http.StatusBadRequest, "Database error.")
				return
			}
		}
		singleVoteResult := models.SingleVoteResult{
			ID:   postId,
			Name: post.Name,
		}
		var candidateResults []models.CandidateResult
		for _, candidate := range post.Candidates {
			tmpBallotId := models.UsedBallotID{
				Name: candidate.Name,
				Roll: candidate.Roll,
			}
			for _, ballotString := range candidate.BallotIDs {
				tmpBallotId.BallotString = ballotString
				ballotIds = append(ballotIds, tmpBallotId)
			}
			candidateResult := models.CandidateResult{
				Name:   candidate.Name,
				Roll:   candidate.Roll,
				Count:  candidate.Count,
				Status: candidate.Status,
			}
			candidateResults = append(candidateResults, candidateResult)

			if candidate.Status == "eliminated" {
				err = ElectionDb.EliminateCandidate(postId, candidate.Roll)
				if err != nil {
					c.String(http.StatusInternalServerError, "Database error.")
					return
				}
			}
		}
		singleVoteResult.Candidates = candidateResults

		err = ElectionDb.InsertSingleVoteResult(&singleVoteResult)
		if err != nil {
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}

		err = utils.ExportBallotIdsToFile(ballotIds, postId)
		if err != nil {
			c.String(http.StatusInternalServerError, "Cannot export BallotIds.")
			return
		}
	}

	Posts, err = ElectionDb.GetPosts()
	if err != nil {
		c.String(http.StatusBadRequest, "Database error.")
		return
	}

	config.ElectionState = config.ResultsAvailable
	c.String(http.StatusAccepted, "Results accepted.")
}
