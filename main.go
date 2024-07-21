package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

/**
 * It initializes the server with an empty map of connections.
 * Returns a pointer to the newly created Server.
 */
func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	log.Println("New connection from client: ", ws.RemoteAddr())
	// It adds the connection to the server's conns map.
	s.conns[ws] = true
	// It then calls the readLoop method to start reading messages from the client.
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		// It continues reading until an error occurs or the connection is closed.
		n, err := ws.Read(buf)
		if err != nil {
			// The connection on the client side has closed
			if err == io.EOF {
				break
			}

			log.Println("Error receiving message: ", err)
			// break (or return) will close the connection
			continue
		}
		// Only the get the amount of buffer used
		msg := buf[:n]
		log.Println("Received message: ", string(msg))
		// It sends a "Hello from server" message back to the client.
		// ws.Write([]byte("Hello from server"))

		// Let everyone know that user sent a message
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	// For each connection, it launches a goroutine to send the message asynchronously.
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				log.Println("Error sending message: ", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	// Handle WebSocket requests on the "/ws" route using the server's handleWS method.
	http.Handle("/ws", websocket.Handler(server.handleWS))
	// Start the server and listen for incoming HTTP requests on port 3000.
	http.ListenAndServe(":3000", nil)
}
