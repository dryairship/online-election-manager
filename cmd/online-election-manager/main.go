package main

import (
    "fmt"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/dryairship/online-election-manager/router"
    "github.com/dryairship/online-election-manager/db"
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/controllers"
)

func main() {
    utils.InitializeRandomSeed()
    
    r := gin.Default()
    
    var err error
    controllers.ElectionDb, err = db.ConnectToDatabase()
    if err != nil {
        fmt.Println("[ERROR] Could not establish database connection.")
        os.Exit(1)
    }
    
    controllers.Posts, err = controllers.ElectionDb.GetPosts()
    if err != nil {
        fmt.Println("[ERROR] Could not set posts data.")
        os.Exit(1)
    }
    
    sessionDb := cookie.NewStore([]byte(config.SessionsKey))
    r.Use(sessions.Sessions("SessionData", sessionDb))
    router.SetUpRoutes(r)
    r.Run(config.ApplicationPort)
}
