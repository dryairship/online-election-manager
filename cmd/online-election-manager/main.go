package main

import (
	"fmt"
	"os"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/controllers"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/router"
	"github.com/dryairship/online-election-manager/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// To seed the random number generator.
	utils.InitializeRandomSeed()
	// To initialize the data from environment variables.
	config.InitializeConfiguration()
	// To authenticate the sender of verification mails
	utils.AuthenticateMailer()

	// Create a default server.
	r := gin.Default()

	var err error

	// Connect to the database.
	controllers.ElectionDb, err = db.ConnectToDatabase()
	if err != nil {
		fmt.Println("[ERROR] Could not establish database connection.")
		os.Exit(1)
	}

	ceo, err := controllers.ElectionDb.GetCEO()
	if err != nil {
		fmt.Println("[WARN] CEO has not registered.")
	} else {
		config.PublicKeyOfCEO = ceo.PublicKey
	}

	// Store the posts in a global variable so there is no need
	// to access the database for the same data repeatedly.
	controllers.Posts, err = controllers.ElectionDb.GetPosts()
	if err != nil {
		fmt.Println("[ERROR] Could not set posts data.")
		os.Exit(1)
	}

	// Set up an encrypted cookie based session storage.
	sessionDb := cookie.NewStore([]byte(config.SessionsKey))
	r.Use(sessions.Sessions("SessionData", sessionDb))

	// Set up the routes and listeners for API calls.
	router.SetUpRoutes(r)

	// Start the server on the specified port.
	r.Run(config.ApplicationPort)
}
