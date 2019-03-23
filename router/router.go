package router

import (
    "github.com/dryairship/online-election-manager/controllers"
    "github.com/dryairship/online-election-manager/config"
    "github.com/gin-gonic/contrib/static"
    "github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
    
    users := r.Group("/users")
    {
        users.GET("/mail/:roll", controllers.SendMailToStudent)
        users.POST("/register", controllers.RegisterNewVoter)
        users.POST("/login", controllers.CheckUserLogin)
    }
    
    election := r.Group("/election")
    {
        election.GET("/getVotablePosts/:roll", controllers.GetVotablePosts)
        election.GET("/getCandidateInfo/:roll", controllers.GetCandidateInfo)
        election.GET("/getCandidateCard/:roll", controllers.GetCandidateCard)
        election.POST("/submitVote", controllers.SubmitVote)
    }
    
    r.Use(static.Serve("/",static.LocalFile(config.AssetsPath,true)))
}
