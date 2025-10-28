package websocket

import (
	"encoding/json"
	"log"
	"project/repository"

	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Hub  *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: unexpected close: %v", err)
			}
			break
		}

		// Basic validation: try to unmarshal into Message
		var msg MessageWs
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("invalid message format", err)
			continue
		}
		//update online status
		c.Hub.UpdateUserOnlineThrottled(c.ID, 10*time.Second)
		switch msg.Type {
		case "ping":
			// Respond to ping
			c.Conn.WriteMessage(websocket.PongMessage, nil)

		case "chat":
			if msg.To != "" {
				c.Hub.SendToUser(msg.To, message)
			}
		case "notify_friend":
			if msg.To != "" {
				c.Hub.NotifyInviteFriend(msg.To, msg.Payload)
			}
		case "is_online":
			// Handle is_online broadcast if needed
			is_online, time_online, err := repository.NewRedisRepository().IsUserOnline(msg.To)
			if err != nil {
				log.Printf("failed to check online status for %s: %v", msg.To, err)
				break
			}
			if is_online && time_online > 0 {
				log.Printf("user %s is online, time online: %v", msg.To, time_online)
			}
			payload := map[string]interface{}{
				"type":        "is_online_response",
				"user_id":     msg.To,
				"is_online":   is_online,
				"time_online": time_online.Seconds(),
			}
			response, err := json.Marshal(payload)
			if err != nil {
				log.Printf("failed to marshal is_online response for %s: %v", msg.To, err)
				break
			}
			log.Println("send to user ", msg.From)
			go c.Hub.SendToUser(msg.From, response)
		default:
			b, _ := json.Marshal(MessageWs{
				Type:    "error",
				Payload: json.RawMessage(`"unknown message type"`),
			})
			c.Send <- b
		}
		// If message has a "To" field, forward to specific user

		// Otherwise broadcast to everyone (example)
		// c.Hub.Broadcast <- message
	}
}
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Send queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
