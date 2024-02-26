package session_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Sessions map[uuid.UUID]*Session

type PlayerScores struct {
	Player1Score uint16
	Player2Score uint16
}

type Session struct {
	SessionId        uuid.UUID
	PlayerScores     map[uuid.UUID]int16
	SessionStartTime time.Time
}

type PlayerUpdate struct {
	Score int16
	X     int16
	Y     int16
}

func (s *Session) GetPlayerScores() map[uuid.UUID]int16 {
	return s.PlayerScores
}

func InitSessionEngine() {
	Sessions = make(map[uuid.UUID]*Session)
}

func HandleUserSession(conn *websocket.Conn, s *Session, userId uuid.UUID) {
	defer conn.Close()

	go clientUpdate(conn, s)

	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message")
			return
		}

		var message *PlayerUpdate
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		fmt.Printf("Received new score: %v\nx: %v\ny: %v\n", message.Score, message.X, message.Y)
		s.PlayerScores[userId] = message.Score

		// Example of writing back a message
		// err = conn.WriteMessage(websocket.TextMessage, msg)
		// if err != nil {
		// 	return
		// }
	}
}

func clientUpdate(conn *websocket.Conn, s *Session) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			sessionJSON, err := json.Marshal(s)
			if err != nil {
				fmt.Println("JSON marshal error:", err)
			}

			err = conn.WriteMessage(websocket.TextMessage, sessionJSON)
			if err != nil {
				fmt.Println("Write message error:", err)
				if err == websocket.ErrCloseSent {
					fmt.Println("Breaking connection")
					done <- true
				}
			}
		}
	}
}
