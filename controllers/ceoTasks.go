package controllers

import (
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

// Struct to accept the keys from the client.
type NewCEOData struct {
    PublicKey       string  `json:"pubkey"`
    PrivateKey      string  `json:"privkey"`
}

// Struct to return candidates' data to the CEO.
type CandidateData struct {
    Name        string
    Username    string
    PrivateKey  string
    PostID      string
}

// API handler to send verification mail to CEO.
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

// API handler to register CEO's account.
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

// API handler to fetch all candidates from the database.
func FetchCandidates(c *gin.Context){
    id, err := utils.GetSessionID(c)
    if err != nil || id != "CEO" {
        c.String(http.StatusForbidden, "Only the CEO can access this.")
        return
    }
    
    candidates, err := ElectionDb.GetAllCandidates()
    if err != nil {
        c.String(http.StatusInternalServerError, "Error while fetching candidates.")
    }
    
    data := make([]CandidateData, len(candidates))
    for i, candidate := range candidates {
        data[i] = CandidateData {
            Name:       candidate.Name,
            Username:   candidate.Username,
            PrivateKey: candidate.PrivateKey,
            PostID:     candidate.PostID,
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
    
    posts, err := ElectionDb.GetPosts()
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

