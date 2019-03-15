package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/contrib/static"
)

func main() {
    r := gin.Default()
    r.POST("/login", func(c *gin.Context) {
        if c.PostForm("roll") == "180561" {
            c.JSON(200, gin.H{
                "success": 1,
                "roll": c.PostForm("roll"),
                "pass": c.PostForm("password"),
            })
        } else {
            c.JSON(200, gin.H{
                "success": 0,
            })
        }
    })
    r.Use(static.Serve("/",static.LocalFile("../../assets",true)))
    r.Run(":9999")
}
