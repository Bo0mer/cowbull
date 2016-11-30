package cowbull

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// ensure remotePlayer satisfies Player.
var _ Player = &remotePlayer{}

type cowsbulls struct {
	Cows  int `json:"cows"`
	Bulls int `json:"bulls"`
}

type digits struct {
	Digits int `json:"digits"`
}

type number struct {
	Number string `json:"number"`
}

type remotePlayer struct {
	c *Client

	mu   sync.RWMutex // guards
	name string

	digits chan digits    // number of digits of the unknown number
	number chan number    // the last guess of the player
	try    chan cowsbulls // the result of the last try to guess the number
}

// RemotePlayer creates a player based on a client.
func RemotePlayer(c *Client) *remotePlayer {
	p := &remotePlayer{
		c:      c,
		digits: make(chan digits),
		number: make(chan number),
		try:    make(chan cowsbulls),
	}

	c.OnMessage("name", func(data string) {
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

	c.OnMessage("think", func(data string) {
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

	c.OnMessage("guess", func(data string) {
		go func() {
			var n number
			err := json.Unmarshal([]byte(data), &n)
			if err != nil {
				log.Printf("remoteplayer: bad input for guess: %s\n", data)
			}
			p.number <- n
		}()
	})

	c.OnMessage("try", func(data string) {
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
func (p *remotePlayer) ID() string {
	return p.c.ID()
}

// Name returns the remote player's name. It may be empty.
func (p *remotePlayer) Name() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.name
}

func (p *remotePlayer) AnnouncePlayers(players []PlayerEntry) error {
	playersBytes, err := json.Marshal(players)
	if err != nil {
		return err
	}
	return p.c.SendMessage("players", string(playersBytes))
}

func (p *remotePlayer) Think() (int, error) {
	if err := p.c.SendMessage("think", ""); err != nil {
		return 0, err
	}

	return (<-p.digits).Digits, nil
}

func (p *remotePlayer) Guess(n int) (string, error) {
	guessReq := digits{Digits: n}
	data, err := json.Marshal(&guessReq)
	if err != nil {
		return "", err
	}

	if err := p.c.SendMessage("guess", string(data)); err != nil {
		return "", err
	}

	return (<-p.number).Number, nil
}

func (p *remotePlayer) Try(guess string) (int, int, error) {
	tryReq := number{Number: guess}
	data, err := json.Marshal(&tryReq)
	if err != nil {
		return 0, 0, err
	}

	if err := p.c.SendMessage("try", string(data)); err != nil {
		return 0, 0, err
	}

	try := <-p.try
	fmt.Printf("returning %#v\n", try)
	return try.Cows, try.Bulls, nil
}

func (p *remotePlayer) Tell(cows, bulls int) error {
	tellReq := cowsbulls{Cows: cows, Bulls: bulls}
	data, err := json.Marshal(&tellReq)
	if err != nil {
		return err
	}
	return p.c.SendMessage("tell", string(data))
}
