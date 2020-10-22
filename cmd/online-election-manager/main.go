package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/controllers"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/router"
	"github.com/dryairship/online-election-manager/utils"
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
		log.Fatalln("[ERROR] Could not establish database connection: ", err.Error())
	}

	ceo, err := controllers.ElectionDb.GetCEO()
	if err != nil {
		log.Println("[WARN] CEO has not registered.")
	} else {
		config.PublicKeyOfCEO = ceo.PublicKey
	}

	// Store the posts in a global variable so there is no need
	// to access the database for the same data repeatedly.
	controllers.Posts, err = controllers.ElectionDb.GetPosts()
	if err != nil {
		log.Fatalln("[ERROR] Could not set posts data:", err.Error())
	}

	// Set up an encrypted cookie based session storage.
	sessionDb := cookie.NewStore([]byte(config.SessionsKey))
	r.Use(sessions.Sessions("SessionData", sessionDb))

	// Set up the routes and listeners for API calls.
	router.SetUpRoutes(r)

	// Start the server on the specified port.
	if err = r.Run(config.ApplicationPort); err != nil {
		log.Fatalln("[ERROR] Could not start the server: ", err.Error())
	}
}
