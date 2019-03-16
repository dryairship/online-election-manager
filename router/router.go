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
    
    r.Use(static.Serve("/",static.LocalFile(config.AssetsPath,true)))
}
