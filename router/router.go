package router

import (
    "github.com/dryairship/online-election-manager/controllers"
    "github.com/dryairship/online-election-manager/config"
    "github.com/gin-gonic/contrib/static"
    "github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
    
    session := r.Group("/session")
    {
        session.POST("/login", controllers.SessionLogin)
    }
    
    users := r.Group("/users")
    {
        users.GET("/mail/:roll", controllers.SendMailToStudent)
        users.POST("/register", controllers.RegisterNewVoter)
    }
    
    r.Use(static.Serve("/",static.LocalFile(config.AssetsPath,false)))
}
