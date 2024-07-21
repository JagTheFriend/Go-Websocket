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

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	log.Println("New connection from client: ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
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
		ws.Write([]byte("Hello from server"))
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":3000", nil)
}
