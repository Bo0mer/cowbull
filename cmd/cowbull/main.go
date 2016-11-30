package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bo0mer/cowbull"
	"github.com/Bo0mer/cowbull/game"
	"github.com/gorilla/websocket"
)

var (
	addr            string
	skipOriginCheck bool
)

const (
	addrUsage        = "Server address."
	checkOriginUsage = "Skip Origin header check upon WebSocket connection negotiation."
)

func init() {
	flag.StringVar(&addr, "address", "127.0.0.1:8080", addrUsage)
	flag.BoolVar(&skipOriginCheck, "skip-origin-check", false, checkOriginUsage)
}

func main() {
	flag.Parse()

	var checkOrigin func(*http.Request) bool
	if skipOriginCheck {
		checkOrigin = func(_ *http.Request) bool {
			return true
		}

	}

	gamer := gamer{}
	playerHub := cowbull.NewHub(gamer, log.New(os.Stdout, "hub: ", 0))
	srv := cowbull.Server(&cowbull.ServerConfig{
		StaticFilesPath: "./static/",
		Log:             log.New(os.Stdout, "server: ", 0),
		Hub:             playerHub,
		Upgrader: &websocket.Upgrader{
			HandshakeTimeout:  time.Second * 5,
			CheckOrigin:       checkOrigin,
			EnableCompression: false,
		},
	})

	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("error serving: %v\n", err)
	}
}

type gamer struct{}

func (gamer) Game(t game.Thinker, g game.Guesser) (*game.Game, error) {
	return game.New(t, g), nil
}
