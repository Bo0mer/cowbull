package cowbull

import (
	"encoding/json"
	"log"
	"sync"
)

//go:generate counterfeiter . Messenger

// ensure RemotePlayer satisfies Player.
var _ Player = &RemotePlayer{}

type cowsbulls struct {
	Number string `json:"number"`
	Cows   int    `json:"cows"`
	Bulls  int    `json:"bulls"`
}

type digits struct {
	Digits int `json:"digits"`
}

type number struct {
	Number string `json:"number"`
}

// Messenger sends, receives and acts on messages.
type Messenger interface {
	// ID is the messenger's id.
	ID() string
	// OnMessage registers an action for a message of a kind.
	OnMessage(kind string, action func(data string))
	// SendMessage sends a message of a kind.
	SendMessage(kind string, data string) error
}

// RemotePlayer is a player that is away and the only way to communicate
// with it is via messenger.
type RemotePlayer struct {
	m Messenger

	mu   sync.RWMutex // guards
	name string

	digits chan digits    // number of digits of the unknown number
	number chan number    // the last guess of the player
	try    chan cowsbulls // the result of the last try to guess the number
}

// NewRemotePlayer creates a player based on a messenger.
func NewRemotePlayer(m Messenger) *RemotePlayer {
	p := &RemotePlayer{
		m:      m,
		digits: make(chan digits),
		number: make(chan number),
		try:    make(chan cowsbulls),
	}

	m.OnMessage("name", func(data string) {
		var name struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal([]byte(data), &name); err != nil {
			log.Printf("remoteplayer: bad input for name: %s\n", data)
			return
		}
		p.mu.Lock()
		defer p.mu.Unlock()
		p.name = name.Name
	})

	m.OnMessage("think", func(data string) {
		go func() {
			var thinkResp digits
			err := json.Unmarshal([]byte(data), &thinkResp)
			if err != nil {
				log.Printf("remoteplayer: bad input for think: %s\n", data)
				return
			}
			p.digits <- thinkResp
		}()
	})

	m.OnMessage("guess", func(data string) {
		go func() {
			var n number
			err := json.Unmarshal([]byte(data), &n)
			if err != nil {
				log.Printf("remoteplayer: bad input for guess: %s\n", data)
			}
			p.number <- n
		}()
	})

	m.OnMessage("try", func(data string) {
		go func() {
			var cb cowsbulls
			err := json.Unmarshal([]byte(data), &cb)
			if err != nil {
				log.Printf("remoteplayer: bad input for try: %s\n", data)
				return
			}
			p.try <- cb
		}()
	})

	return p
}

// ID returns the remote player's id.
func (p *RemotePlayer) ID() string {
	return p.m.ID()
}

// Name returns the remote player's name. It may be empty.
func (p *RemotePlayer) Name() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.name
}

// AnnouncePlayers sends message announcing all players specified.
func (p *RemotePlayer) AnnouncePlayers(players []PlayerEntry) error {
	playersBytes, err := json.Marshal(players)
	if err != nil {
		return err
	}
	return p.m.SendMessage("players", string(playersBytes))
}

// Think sends a think messages and returns its response.
func (p *RemotePlayer) Think() (int, error) {
	if err := p.m.SendMessage("think", ""); err != nil {
		return 0, err
	}

	return (<-p.digits).Digits, nil
}

// Guess sends a guess message and returns its response.
func (p *RemotePlayer) Guess(n int) (string, error) {
	guessReq := digits{Digits: n}
	data, err := json.Marshal(&guessReq)
	if err != nil {
		return "", err
	}

	if err := p.m.SendMessage("guess", string(data)); err != nil {
		return "", err
	}

	return (<-p.number).Number, nil
}

// Try sends a try message and returns its response.
func (p *RemotePlayer) Try(guess string) (int, int, error) {
	tryReq := number{Number: guess}
	data, err := json.Marshal(&tryReq)
	if err != nil {
		return 0, 0, err
	}

	if err := p.m.SendMessage("try", string(data)); err != nil {
		return 0, 0, err
	}

	try := <-p.try
	return try.Cows, try.Bulls, nil
}

// Tell sends a tell message.
func (p *RemotePlayer) Tell(number string, cows, bulls int) error {
	tellReq := cowsbulls{Number: number, Cows: cows, Bulls: bulls}
	data, err := json.Marshal(&tellReq)
	if err != nil {
		return err
	}
	return p.m.SendMessage("tell", string(data))
}
