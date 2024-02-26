package routes

import (
	"fmt"
	"net/http"
	"time"

	"example.com/server/session_manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections | TODO Make specific for web app
		return true
	},
}

func GetSession(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	session := session_manager.Sessions[uuid]
	c.JSON(http.StatusOK, session)
}

func SetUserId(c *gin.Context) {
	userId := uuid.New().String()
	c.SetCookie("userId", userId, 3600 /* age */, "/" /* valid for all paths */, "localhost", false /* HTTPS only */, false)
	c.JSON(http.StatusOK, userId)
}

func ConnectSession(c *gin.Context) {
	var session *session_manager.Session = nil
	var userId uuid.UUID
	fmt.Println("Connecting...")
	// Ensure the client has a userId cookie set
	if cookieVal, err := c.Cookie("userId"); err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusNotFound, gin.H{"error": "No userId cookie set, cannot upgrade connection"})
			return
		}

	} else {
		uuid, err := uuid.Parse(cookieVal)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Could not parse cookie into a valid uuid"})
			return
		} else {
			userId = uuid
		}
	}

	// Verify a session id was passed in
	id := c.Query("id")
	fmt.Println("Id is: " + id)
	sessionId, err := uuid.Parse(id)
	if err != nil {
		sessionId = uuid.New()
		session = &session_manager.Session{SessionId: sessionId, SessionStartTime: time.Now()}
		session.PlayerScores = make(map[uuid.UUID]int16)
		session_manager.Sessions[session.SessionId] = session
		session.PlayerScores[userId] = 0
	} else {
		session_manager.Sessions[sessionId].PlayerScores[userId] = 0
	}

	// Verify that the session was successfully created
	session, exists := session_manager.Sessions[sessionId]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Session '%s' does not exist", sessionId)})
		return
	}

	// Upgrade HTTP connection to WebSocket
	fmt.Println("Attempting to upgrade to websocket connection...")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Failed to upgrade to websocket connection", http.StatusInternalServerError)
		return
	}

	// Handle WebSocket communication here
	go session_manager.HandleUserSession(conn, session, userId)
}

func CreateSession(c *gin.Context) {
	uuid := uuid.New()
	s := &session_manager.Session{SessionId: uuid, SessionStartTime: time.Now()}

	session_manager.Sessions[uuid] = s
	// Send a response
	c.JSON(http.StatusOK, s)
}
