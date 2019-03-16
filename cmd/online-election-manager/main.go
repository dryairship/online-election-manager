package main

import (
    "fmt"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/dryairship/online-election-manager/utils"
    "github.com/dryairship/online-election-manager/router"
    "github.com/dryairship/online-election-manager/db"
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/controllers"
)

func main() {
    r := gin.Default()
    router.SetUpRoutes(r)
    utils.InitializeRandomSeed()
    var err error
    controllers.ElectionDb, err = db.ConnectToDatabase()
    if err != nil {
        fmt.Println("[ERROR] Could not establish database connection.")
        os.Exit(1)
    } else {
        r.Run(config.ApplicationPort)
    }
}
