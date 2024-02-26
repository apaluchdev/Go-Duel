package main

import (
	"example.com/server/routes"
	"example.com/server/session_manager"
	"example.com/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	reactAppDomain := "http://localhost:3000" //os.Getenv("REACT_APP_DOMAIN")

	// Allow requests from your React app's origin
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", reactAppDomain)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	session_manager.InitSessionEngine()
	sessionGroup := r.Group("/session")
	{
		sessionGroup.GET("/:session", routes.CreateSession)
		sessionGroup.POST("/:session", routes.GetSession)
		sessionGroup.GET("/setuserid", routes.SetUserId)
		sessionGroup.GET("/connect", routes.ConnectSession)
	}

	r.Static("/static", "../static")

	r.Run()

	// Testing how separate modules work
	utils.Greet("Bob")
}
