package utils

import (
	"errors"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// This function stores the identification of the logged in user
// in the encrypted cookie-based session storage.
func StartSession(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("UserType") != nil {
		session.Clear()
		if err := session.Save(); err != nil {
			log.Println("[ERROR] Could not save empty session")
		}
	}

	roll := c.PostForm("roll")
	session.Set("UserType", "Authenticated")
	session.Set("ID", roll)
	if err := session.Save(); err != nil {
		log.Println("[ERROR] Could not save session for ", roll)
	}
}

// This function deletes the data store about the user from
// the session storage.
func EndSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Println("[ERROR] Could not end session")
	}
}

// This function returns the identification of the logged in user.
func GetSessionID(c *gin.Context) (string, error) {
	id := sessions.Default(c).Get("ID")
	if id != nil {
		return id.(string), nil
	}
	return "", errors.New("Unauthenticated User")
}
