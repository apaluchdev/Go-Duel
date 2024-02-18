package main

import (
	"example.com/server/routes"
	"example.com/server/session_manager"
	"example.com/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Instruct browsers to not cache
	r.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	})

	session_manager.InitSessionEngine()
	sessionGroup := r.Group("/session")
	{
		sessionGroup.GET("/:session", routes.CreateSession)
		sessionGroup.POST("/:session", routes.GetSession)
		sessionGroup.GET("/setuserid", routes.SetUserId)
		sessionGroup.GET("/connect", routes.ConnectSession)
		sessionGroup.GET("/ws", routes.WebSocketExample) // Example websocket connection
	}

	r.Static("/static", "../static")

	r.Run()

	// Testing how separate modules work
	utils.Greet("Bob")
}
