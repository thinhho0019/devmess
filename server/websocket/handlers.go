package websocket

import (
	"log"
	"net/http"
	"project/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	Hub         *Hub
	AuthService *service.AuthService
}
func NewWsHandler(hub *Hub, authService *service.AuthService) *WsHandler {
	return &WsHandler{
		Hub:         hub,
		AuthService: authService,
	}
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("🔍 [WebSocket] Connection from origin: %s", origin)
		return true
	},
}

func (h *WsHandler) ServeWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("🔌 [ServeWs] WebSocket connection attempt")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("❌ [ServeWs] Upgrade error: %v", err)
			return
		}

		log.Println("✅ [ServeWs] WebSocket upgrade successful")

		token := c.Query("token")
		log.Printf("🔑 [ServeWs] Received token: %s", token)

		if token == "" {
			log.Println("❌ [ServeWs] Missing token in query parameters")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Missing token"))
			conn.Close()
			return
		}

		log.Printf("🔐 [ServeWs] Verifying token for WebSocket connection: %s", token)
		user, _, err := h.AuthService.VerifyAccessToken(token)
		if err != nil {
			log.Printf("❌ [ServeWs] Authentication failed: %v", err)
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Authentication failed"))
			conn.Close()
			return
		}
		log.Printf("✅ [ServeWs] User %s authenticated and connected", user.ID.String())

		client := &Client{
			ID:   user.ID.String(),
			Hub:  h.Hub,
			Conn: conn,
			Send: make(chan []byte, 256),
		}

		log.Printf("📝 [ServeWs] Registering client %s with hub", user.ID.String())
		client.Hub.register <- client

		log.Printf("🚀 [ServeWs] Starting pumps for client %s", user.ID.String())

		// FIX: Start ReadPump as goroutine, WritePump blocks
		go client.ReadPump()

		// WritePump MUST block to keep the HTTP handler alive
		// When WritePump returns, the connection is closed
		client.WritePump()

		log.Printf("🔌 [ServeWs] Connection closed for user %s", user.ID.String())
	}
}
