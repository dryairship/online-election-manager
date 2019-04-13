package utils

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "errors"
)

// This function stores the identification of the logged in user
// in the encrypted cookie-based session storage.
func StartSession(c *gin.Context) {
    session := sessions.Default(c)
    
    if session.Get("UserType") != nil {
        session.Clear()
        session.Save()
    }
    
    roll := c.PostForm("roll")
    session.Set("UserType","Authenticated")
    session.Set("ID",roll)
    session.Save()
}

// This function deletes the data store about the user from
// the session storage.
func EndSession(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Save()
}

// This function returns the identification of the logged in user.
func GetSessionID(c *gin.Context) (string, error) {
    id := sessions.Default(c).Get("ID")
    if id != nil {
        return id.(string), nil
    }
    return "", errors.New("Unauthenticated User")
}
