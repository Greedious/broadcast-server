package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	clients   map[*websocket.Conn]string
	broadcast chan []byte
	mu        sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan []byte),
	}
}

func (s *Server) DisconnectClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer conn.Close()
	delete(s.clients, conn)
	s.mu.Unlock()
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	// Get client name from query
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Anonymous"
	}

	fmt.Printf("%s client has joined the chat\n", name)

	s.mu.Lock()
	s.clients[conn] = name
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		fmt.Printf("Client '%s' disconnected\n", s.clients[conn])
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Prepend name to the message
		formattedMsg := fmt.Sprintf("%s: %s", name, string(msg))
		s.broadcast <- []byte(formattedMsg)
	}
}
func (s *Server) HandleMessages() {
	for {
		msg := <-s.broadcast
		if msg == nil {
			continue
		}

		s.mu.Lock()
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mu.Unlock()
	}
}
