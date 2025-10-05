package handler

import (
	"net/http"
	"project/websocket"
)

func ServeWs(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
	websocket.ServeWs(hub, w, r)
}
