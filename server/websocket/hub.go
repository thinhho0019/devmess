// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"errors"
	"fmt"
	"log"
	"project/repository"
	"sync"
	"time"
)

// Message represents a message sent to clients.

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clientsMu sync.RWMutex
	clients   map[string]*Client // userID -> client

	// Registered requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Inbound messages from the clients.
	Broadcast chan []byte

	lastUpdateMu sync.RWMutex
	lastUpdate   map[string]time.Time // userID -> last update timestamp
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		lastUpdate: make(map[string]time.Time),
	}
}

func (h *Hub) Run() {
	cleanupTicker := time.NewTicker(5 * time.Minute)
	defer cleanupTicker.Stop()
	for {
		select {
		case client := <-h.register:
			h.clientsMu.Lock()
			h.clients[client.ID] = client
			h.clientsMu.Unlock()
			log.Printf("client registered: %s", client.ID)
			// Update user online status asynchronously if repository exposes such method.
			go func(user_id string) {
				// If repository.RedisRepository provides UpdateUserOnlineStatus, uncomment and use it:
				err := repository.NewRedisRepository().UpdateUserOnline(user_id)
				if err != nil {
					log.Printf("failed to update user online status for %s: %v", user_id, err)
				}
				log.Printf("updateUserOnlineStatus not called for: %s (method not found)", user_id)
			}(client.ID)

		case client := <-h.unregister:
			h.clientsMu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.clientsMu.Unlock()
			log.Printf("client unregistered: %s", client.ID)

		case message := <-h.Broadcast:
			// h.clientsMu.RLock()
			// for _, c := range h.clients {
			// 	select {
			// 	case c.Send <- message:
			// 	default:
			// 		// slow client, drop
			// 		close(c.Send)
			// 		delete(h.clients, c.ID)
			// 	}
			// }
			// h.clientsMu.RUnlock()
			println(message)
			continue
		case <-cleanupTicker.C:
            h.cleanupLastUpdate()
        
		}
	}
}
func (h *Hub) cleanupLastUpdate() {

	h.lastUpdateMu.Lock()
	h.clientsMu.RLock()

	// Remove entries for disconnected users
	for uid := range h.lastUpdate {
		if _, connected := h.clients[uid]; !connected {
			delete(h.lastUpdate, uid)
		}
	}

	h.clientsMu.RUnlock()
	h.lastUpdateMu.Unlock()

	log.Printf("🧹 [Hub] Cleaned up stale lastUpdate entries")
}

func (h *Hub) SendToUser(userID string, message []byte) error {
	print("Sending message to user:", userID)
	h.clientsMu.RLock()
	c, ok := h.clients[userID]
	h.clientsMu.RUnlock()
	if !ok {
		return errors.New("user not connected")
	}

	select {
	case c.Send <- message:
		return nil
	default:
		return errors.New("user send channel full")
	}
}
func (h *Hub) SendToUsers(users []string, message []byte) error {
	var failed []string
	for _, userID := range users {
		h.clientsMu.RLock()
		c, ok := h.clients[userID]
		h.clientsMu.RUnlock()
		if !ok {
			failed = append(failed, userID)
			continue
		}
		select {
		case c.Send <- message:
			println("send message success")
		default:
			failed = append(failed, userID)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to send to users: %v", failed)
	}
	return nil
}

func (h *Hub) NotifyInviteFriend(userID string, message []byte) error {
	h.clientsMu.RLock()
	c, ok := h.clients[userID]
	h.clientsMu.RUnlock()
	if !ok {
		return errors.New("user not connected")
	}
	select {
	case c.Send <- message:
		return nil
	default:
		return errors.New("user send channel full")
	}
}

func (h *Hub) UpdateUserOnlineThrottled(userID string, threshold time.Duration) {
	h.lastUpdateMu.Lock()
	lastTime, exists := h.lastUpdate[userID]
	now := time.Now()

	// Skip if updated recently
	if exists && now.Sub(lastTime) < threshold {
		h.lastUpdateMu.Unlock()
		return
	}

	// Update timestamp
	h.lastUpdate[userID] = now
	h.lastUpdateMu.Unlock()

	// Update Redis (non-blocking)
	go func(uid string) {
		if err := repository.NewRedisRepository().UpdateUserOnline(uid); err != nil {
			log.Printf("⚠️ [Hub] Failed to update online for %s: %v", uid, err)
		} else {
			log.Printf("🔄 [Hub] Online status updated for %s", uid)
		}
	}(userID)
}
