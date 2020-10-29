package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

// Struct to accept the keys from the client.
type newCEOData struct {
	PublicKey  string `json:"pubkey"`
	PrivateKey string `json:"privkey"`
}

// Struct to return candidates' data to the CEO.
type candidateData struct {
	Name       string
	Roll       string
	PostID     string
	PrivateKey string
}

type singleVoteCandidateResult struct {
	Roll      string   `json:"roll"`
	Name      string   `json:"name"`
	Count     int32    `json:"count"`
	BallotIDs []string `json:"ballotIds"`
	Status    string   `json:"status"`
}

type candidateResult struct {
	Roll        string `json:"roll"`
	Name        string `json:"name"`
	Preference1 int32  `json:"preference1"`
	Preference2 int32  `json:"preference2"`
	Preference3 int32  `json:"preference3"`
}

type singleVotePostResult struct {
	ID         string                      `json:"postid"`
	Name       string                      `json:"postname"`
	Resolved   bool                        `json:"resolved"`
	Candidates []singleVoteCandidateResult `json:"candidates"`
}

type postResult struct {
	ID         string            `json:"postid"`
	Name       string            `json:"postname"`
	Candidates []candidateResult `json:"candidates"`
	BallotIDs  []string          `json:"ballotIds"`
}

type singleVoteResultData struct {
	Posts []singleVotePostResult `json:"posts"`
}

type resultData struct {
	Posts []postResult `json:"posts"`
}

// API handler to send verification mail to CEO.
func SendMailToCEO(c *gin.Context) {
	ceo, err := ElectionDb.GetCEO()

	// Checking for err == nil because ceo is inserted in the database when
	// the verification mail is sent to them,
	if err == nil {
		log.Println("[WARN] CEO re-requested verification mail")
		c.String(http.StatusForbidden, "Verification mail has already been sent to CEO.")
		return
	}

	skeleton, err := ElectionDb.FindStudentSkeleton(config.RollNumberOfCEO)
	if err != nil {
		log.Println("[ERROR] No CEO in the database, can't send verification mail")
		c.String(http.StatusInternalServerError, "No CEO assigned.")
		return
	}

	ceo = skeleton.CreateCEO(utils.GetRandomAuthCode())
	recipient := ceo.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "the CEO")
	if err != nil {
		log.Println("[ERROR] Database error while sending mail to CEO: ", ceo, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	} else {
		err = ElectionDb.InsertCEO(&ceo)
		if err != nil {
			log.Println("[ERROR] Database error while inserting CEO: ", ceo, err.Error())
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
		log.Println("[WARN] CEO tried to re-register")
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
		log.Println("[ERROR] Database error while registering CEO: ", ceo, err.Error())
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
		log.Println("[ERROR] CEO Login attempted, but CEO has not registered")
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
		log.Println("[WARN] Unauthorized fetchVotes attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	if config.ElectionState != config.VotingStopped {
		log.Println("[WARN] CEO tried to fetch votes before stopping voting")
		c.String(http.StatusForbidden, "Please stop voting before fetching votes")
		return
	}

	votes, err := ElectionDb.GetVotes()
	if err != nil {
		log.Println("[ERROR] Database error while fetching votes: ", err.Error())
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
		log.Println("[WARN] Unauthorized fetchCandidates attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	candidates, err := ElectionDb.GetAllCandidates()
	if err != nil {
		log.Println("[ERROR] Database error while fetching candidates: ", err.Error())
		c.String(http.StatusInternalServerError, "Error while fetching candidates.")
	}

	var data []candidateData
	for _, candidate := range candidates {
		if isCandidateStillInRace(candidate.PostID, candidate.Username) {
			data = append(data, candidateData{
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
		log.Println("[WARN] Unauthorized fetchPosts attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	posts, err := ElectionDb.GetPostsForCEO()
	if err != nil {
		log.Println("[ERROR] Database error while fetching posts: ", err.Error())
		c.String(http.StatusInternalServerError, "Error while fetching posts.")
		return
	}
	c.JSON(http.StatusOK, &posts)
}

// API handler to start voting process by accepting CEO's public and private keys.
func StartVoting(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		log.Println("[WARN] Unauthorized startVoting attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	if config.ElectionState != config.VotingNotYetStarted {
		log.Println("[WARN] CEO tried to start voting, but voting has already started")
		c.String(http.StatusBadRequest, "You are too late. Voting has already started.")
		return
	}

	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		log.Println("[ERROR] No CEO in the database error on starting voting: ", err.Error())
		c.String(http.StatusInternalServerError, "CEO not yet assigned.")
		return
	}
	newData := newCEOData{}
	err = c.BindJSON(&newData)
	if err != nil {
		log.Println("[ERROR] CEO data JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	ceo.PublicKey = newData.PublicKey
	ceo.PrivateKey = newData.PrivateKey
	err = ElectionDb.UpdateCEO(&ceo)
	if err != nil {
		log.Println("[ERROR] Database error while starting voting: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	config.PublicKeyOfCEO = ceo.PublicKey
	config.ElectionState = config.AcceptingVotes

	log.Println("[INFO] Voting started")
	c.String(http.StatusOK, "Voting Started.")
}

// API handler to stop voting process.
func StopVoting(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		log.Println("[WARN] Unauthorized stopVoting attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	if config.ElectionState != config.AcceptingVotes {
		log.Println("[WARN] CEO tried to stop voting but current state = ", config.ElectionState)
		c.String(http.StatusBadRequest, "The system is already not accepting new votes.")
		return
	}

	config.ElectionState = config.VotingStopped
	log.Println("[INFO] Voting stopped")
	c.String(http.StatusOK, "Voting Stopped")
}

func ResultProgress(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("%.3f%%", config.ResultProgress))
}

func PrepareForNextRound(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		log.Println("[WARN] Unauthorized prepareForNextRound attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	err = ElectionDb.MarkAllVotersUnvoted()
	if err != nil {
		log.Println("[ERROR] Database error while marking all voters unvoted: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	err = ElectionDb.DeleteAllVotes()
	if err != nil {
		log.Println("[ERROR] Database error while deleting all votes: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	config.ElectionState = config.VotingNotYetStarted

	log.Println("[INFO] Ready for next round")
	c.String(http.StatusOK, "Ready for next round.")
}

func SubmitSingleVoteResults(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		log.Println("[WARN] Unauthorized submitSingleVoteResults attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	var results singleVoteResultData
	err = c.BindJSON(&results)
	if err != nil {
		log.Println("[ERROR] SingleVoteResults JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, post := range results.Posts {
		var ballotIds []models.UsedSingleVoteBallotID
		postId := post.ID
		if post.Resolved {
			if err := ElectionDb.MarkPostResolved(postId); err != nil {
				log.Println("[ERROR] Database error while marking post resolved: ", postId, err.Error())
				c.String(http.StatusBadRequest, "Database error.")
				return
			}
		}
		singleVoteResult := models.SingleVoteResult{
			ID:   postId,
			Name: post.Name,
		}
		var candidateResults []models.SingleVoteCandidateResult
		for _, candidate := range post.Candidates {
			tmpBallotId := models.UsedSingleVoteBallotID{
				Name: candidate.Name,
				Roll: candidate.Roll,
			}
			for _, ballotString := range candidate.BallotIDs {
				tmpBallotId.BallotString = ballotString
				ballotIds = append(ballotIds, tmpBallotId)
			}
			candidateResult := models.SingleVoteCandidateResult{
				Name:   candidate.Name,
				Roll:   candidate.Roll,
				Count:  candidate.Count,
				Status: candidate.Status,
			}
			candidateResults = append(candidateResults, candidateResult)

			if candidate.Status == "eliminated" {
				err = ElectionDb.EliminateCandidate(postId, candidate.Roll)
				if err != nil {
					log.Println("[ERROR] Database error while eliminating candidate: ", postId, candidate, err.Error())
					c.String(http.StatusInternalServerError, "Database error.")
					return
				}
			}
		}
		singleVoteResult.Candidates = candidateResults

		err = ElectionDb.InsertSingleVoteResult(&singleVoteResult)
		if err != nil {
			log.Println("[ERROR] Database error while inserting SingleVoteResult: ", singleVoteResult, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}

		err = utils.ExportSingleVoteBallotIdsToFile(ballotIds, postId)
		if err != nil {
			log.Println("[ERROR] Util error while exporting ballot IDs to file: ", err.Error())
			c.String(http.StatusInternalServerError, "Cannot export BallotIds.")
			return
		}
	}

	Posts, err = ElectionDb.GetPosts()
	if err != nil {
		log.Println("[ERROR] Database error while getting remaining posts: ", err.Error())
		c.String(http.StatusBadRequest, "Database error.")
		return
	}

	config.ElectionState = config.ResultsAvailable
	log.Println("[INFO] Results stored in the database")
	c.String(http.StatusAccepted, "Results accepted.")
}

func SubmitResults(c *gin.Context) {
	id, err := utils.GetSessionID(c)
	if err != nil || id != "CEO" {
		log.Println("[WARN] Unauthorized submitResults attempt: ", id, err.Error())
		c.String(http.StatusForbidden, "Only the CEO can access this.")
		return
	}

	var results resultData
	err = c.BindJSON(&results)
	if err != nil {
		log.Println("[ERROR] Results JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, post := range results.Posts {
		postId := post.ID
		result := models.Result{
			ID:   postId,
			Name: post.Name,
		}
		var candidateResults []models.CandidateResult
		for _, candidate := range post.Candidates {
			candidateResult := models.CandidateResult{
				Name:        candidate.Name,
				Roll:        candidate.Roll,
				Preference1: candidate.Preference1,
				Preference2: candidate.Preference2,
				Preference3: candidate.Preference3,
				Status:      "none",
			}
			candidateResults = append(candidateResults, candidateResult)
		}
		result.Candidates = candidateResults

		err = ElectionDb.InsertResult(&result)
		if err != nil {
			log.Println("[ERROR] Database error while inserting Result: ", result, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}

		err = utils.ExportBallotIdsToFile(post.BallotIDs, postId)
		if err != nil {
			log.Println("[ERROR] Util error while exporting ballot IDs to file: ", err.Error())
			c.String(http.StatusInternalServerError, "Cannot export BallotIds.")
			return
		}
	}

	Posts, err = ElectionDb.GetPosts()
	if err != nil {
		log.Println("[ERROR] Database error while getting remaining posts: ", err.Error())
		c.String(http.StatusBadRequest, "Database error.")
		return
	}

	config.ElectionState = config.ResultsAvailable
	log.Println("[INFO] Results stored in the database")
	c.String(http.StatusAccepted, "Results accepted.")
}
