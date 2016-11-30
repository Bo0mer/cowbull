package cowbull

import (
	"errors"
	"log"

	"github.com/Bo0mer/cowbull/game"
)

// PlayerEntry holds metadata for a player.
type PlayerEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Player represents a hub member.
type Player interface {

	// ID returns the player's id. It must be unique.
	ID() string

	// Name returns the player's name.
	Name() string

	// AnnouncePlayers announces all players joined the hub.
	AnnouncePlayers([]PlayerEntry) error

	game.Thinker
	game.Guesser
}

// GameSettings represents settings for a game request to a Hub.
type GameSettings struct {
	Role      string   `json:"role"`      // role of the player, either thinker or guesser
	Digits    int      `json:"digits"`    // how many digits should the number have
	AI        bool     `json:"ai"`        // whether the game is versus AI
	Opponents []string `json:"opponents"` // the opponents of the player
}

type hubOp func(map[string]Player)

// Hub represents a group of players.
type Hub struct {
	players map[string]Player

	log *log.Logger

	ops chan hubOp
}

// NewHub creates a brand new hub.
func NewHub(log *log.Logger) *Hub {
	hub := &Hub{
		players: make(map[string]Player),
		log:     log,
		ops:     make(chan hubOp, 1),
	}
	go hub.loop()

	return hub
}

// Add adds a player to the hub.
// Once added, it will get updates by the hub for any significant events.
func (h *Hub) Add(p Player) {
	h.ops <- func(players map[string]Player) {
		players[p.ID()] = p
		h.log.Printf("player %s joined", p.ID())
	}
	h.broadcastPlayers()
}

// Remove removes a player from the hub.
func (h *Hub) Remove(pid string) {
	h.ops <- func(players map[string]Player) {
		delete(players, pid)
		h.log.Printf("player %s left", pid)
	}
	h.broadcastPlayers()
}

// Players returns all player entries in the hub.
func (h *Hub) Players() []PlayerEntry {
	playersChan := make(chan []PlayerEntry, 1)
	h.ops <- func(players map[string]Player) {
		i := 0
		ret := make([]PlayerEntry, len(players))

		for _, p := range players {
			ret[i].ID = p.ID()
			ret[i].Name = p.Name()
			i++
		}
		playersChan <- ret
	}
	return <-playersChan
}

// NewGame creates a new Game based on the provided settings.
func (h *Hub) NewGame(from Player, settings GameSettings) (*game.Game, error) {
	if settings.AI {
		return h.newGameWithAI(from, settings.Digits)
	}
	opponents := h.playersWithIDs(settings.Opponents)
	if len(opponents) == 0 {
		return nil, errors.New("no opponents found")
	}

	if settings.Role == "thinker" {
		return game.New(from, opponents[0]), nil
	}
	return game.New(opponents[0], from), nil
}

func (h *Hub) newGameWithAI(from Player, digits int) (*game.Game, error) {
	return game.WithAI(from, digits), nil
}

func (h *Hub) broadcastPlayers() {
	playerIDs := h.Players()
	h.ops <- func(players map[string]Player) {
		for _, p := range players {
			if err := p.AnnouncePlayers(playerIDs); err != nil {
				h.log.Printf("error announcing players to %s: %v", p.ID(), err)
			}
		}
	}
}

func (h *Hub) playersWithIDs(ids []string) []Player {
	playersChan := make(chan []Player, 1)
	h.ops <- func(players map[string]Player) {
		var ret []Player
		for _, id := range ids {
			if p, ok := players[id]; ok {
				ret = append(ret, p)
			}
		}
		playersChan <- ret
	}
	return <-playersChan
}

func (h *Hub) loop() {
	for op := range h.ops {
		op(h.players)
	}
}
