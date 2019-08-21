package router

import (
	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/controllers"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Groups various API calls and routes them to the
// respective controller functions.
func SetUpRoutes(r *gin.Engine) {

	users := r.Group("/users")
	{
		users.POST("/mail/:roll", controllers.SendMailToStudent)
		users.POST("/register", controllers.RegisterNewVoter)
		users.POST("/login", controllers.CheckUserLogin)
		users.GET("/captcha", controllers.GetCaptcha)
	}

	election := r.Group("/election")
	{
		election.GET("/getVotablePosts", controllers.GetVotablePosts)
		election.GET("/getCandidateCard/:username", controllers.GetCandidateCard)
		election.GET("/getElectionState", controllers.GetElectionState)
		election.POST("/submitVote", controllers.SubmitVote)
	}

	ceo := r.Group("/ceo")
	{
		ceo.POST("/startVoting", controllers.StartVoting)
		ceo.POST("/stopVoting", controllers.StopVoting)
		ceo.GET("/fetchPosts", controllers.FetchPosts)
		ceo.GET("/fetchVotes", controllers.FetchVotes)
		ceo.GET("/fetchCandidates", controllers.FetchCandidates)
		ceo.GET("/calculateResult", controllers.CalculateResult)
		ceo.GET("/resultProgress", controllers.ResultProgress)
	}

	candidate := r.Group("/candidate")
	{
		candidate.POST("/confirmCandidature", controllers.ConfirmCandidature)
		candidate.POST("/declarePrivateKey", controllers.DeclarePrivateKey)
	}

	// To directly serve static files in the AssetsPath directory.
	r.Use(static.Serve("/", static.LocalFile(config.AssetsPath, true)))
}
