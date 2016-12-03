package cowbull

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// ServerConfig holds configuration for a server.
type ServerConfig struct {
	// Path to directory from where to serve content.
	StaticFilesPath string

	// Log is used to log messages from the server and its clients.
	Log *log.Logger

	// Upgrader used for upgrading from an HTTP connection to a WebSocket
	// connection.
	Upgrader *websocket.Upgrader

	// Hub for connected players.
	Hub *Hub
}

// Server implements a cowbull game server.
type Server struct {
	mux     http.Handler
	websock *websocket.Upgrader
	log     *log.Logger

	hub *Hub

	fs http.Handler
}

// NewServer creates a new server.
// After once initialized the config object should not be modified.
func NewServer(cfg *ServerConfig) *Server {
	fs := http.FileServer(http.Dir(cfg.StaticFilesPath))
	mux := http.NewServeMux()

	s := &Server{
		mux:     mux,
		websock: cfg.Upgrader,
		log:     cfg.Log,
		fs:      fs,
		hub:     cfg.Hub,
	}

	mux.Handle("/", s.fs)
	mux.HandleFunc("/websocket", s.upgrade)

	return s
}

// ServeHTTP serves HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}

// upgrade upgrades an HTTP connection to a WebSocket connection and
// forks off a client of the WebSocket connection.
func (s *Server) upgrade(w http.ResponseWriter, req *http.Request) {
	conn, err := s.websock.Upgrade(w, req, nil)
	if err != nil {
		s.log.Printf("error upgrading to websocket conn: %v\n", err)
		return
	}

	c := NewClient(conn)
	player := NewRemotePlayer(c, 60*time.Second)
	c.OnMessage("connect", func(_ string) {
		s.log.Printf("client connected: %s\n", c.ID())
		s.hub.Add(player)
	})

	c.OnMessage("disconnect", func(_ string) {
		s.log.Printf("client disconnected: %s\n", c.ID())
		s.hub.Remove(player.ID())
	})

	c.OnMessage("play", func(data string) {
		s.log.Printf("game initiated by player %s with settings %s\n", player.ID(), data)

		var settings GameSettings
		if err := json.Unmarshal([]byte(data), &settings); err != nil {
			s.log.Printf("malformed settings provided by %s\n", player.ID())
			return
		}

		go func() {
			game, err := s.hub.NewGame(player, settings)
			if err != nil {
				s.log.Printf("error creating game: %v\n", err)
				return
			}
			if err := game.Play(); err != nil {
				s.log.Printf("error running game: %v\n", err)
				return
			}
			s.log.Printf("game finished\n")
		}()
	})
}
