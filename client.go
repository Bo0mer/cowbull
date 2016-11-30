package cowbull

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
)

//go:generate counterfeiter . Conn

type message struct {
	// Name identifies a message family.
	Name string `json:"name"`
	// Data is the actual data of the message.
	Data string `json:"data"`
}

// Conn represents a client connection.
type Conn interface {
	// SetReadDeadline sets the read deadline on the underlying
	// network connection.
	SetReadDeadline(t time.Time) error
	// WriteJSON writes the JSON encoding of v to the connection.
	WriteJSON(v interface{}) error
	// ReadJSON reads the next JSON-encoded message from the connection
	// and stores it in the value pointed to by v.
	ReadJSON(v interface{}) error
	// RemoteAddr should return the remote network address.
	RemoteAddr() net.Addr
	// Close should close the
	Close() error
}

// Client is a remote client.
type Client struct {
	id string

	conn          Conn
	readTimeout   time.Duration
	retryCount    int
	retryInterval time.Duration

	mu      sync.Mutex // guards actions
	actions map[string]func(data string)

	log *log.Logger

	closeOnce sync.Once
}

// ClientOption configures a client.
type ClientOption func(c *Client)

// RetryCount configures the number of attempts when trying to read from a conn.
// If read fails more than n times, the conn is considered broken.
// Defaults to 3.
func RetryCount(n int) ClientOption {
	return func(c *Client) {
		c.retryCount = n
	}
}

// RetryInterval configures the interval between consecutive attempts to read
// from a conn.
// Defaults to 1 second.
func RetryInterval(d time.Duration) ClientOption {
	return func(c *Client) {
		c.retryInterval = d
	}
}

// ReadTimeout configures the ReadTimeout for the underlying connection.
func ReadTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.readTimeout = d
	}
}

// LogTo configures the logging destination for a client.
func LogTo(log *log.Logger) ClientOption {
	return func(c *Client) {
		c.log = log
	}
}

// NewClient creates a new client based on conn and applies all options to it.
// After finished using a client, one should call it's Close method.
func NewClient(conn Conn, opts ...ClientOption) *Client {
	rawid := fmt.Sprintf("%s%d", conn.RemoteAddr().String(), time.Now().UnixNano())
	sha := sha256.Sum256([]byte(rawid))
	id := base64.URLEncoding.EncodeToString(sha[:])
	c := &Client{
		id:            id,
		conn:          conn,
		readTimeout:   120 * time.Second,
		retryCount:    3,
		retryInterval: 1 * time.Second,
		actions:       make(map[string]func(string)),
	}
	for _, op := range opts {
		op(c)
	}
	if c.log == nil {
		c.log = log.New(ioutil.Discard, "", 0)
	}
	go c.readLoop()
	return c
}

// ID returns the id of the client.
func (c *Client) ID() string {
	return c.id
}

// OnMessage registers an action for handling specific messages.
func (c *Client) OnMessage(name string, action func(data string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.actions[name] = action
}

// SendMessage sends message to the client.
func (c *Client) SendMessage(name, data string) error {
	return c.conn.WriteJSON(&message{
		Name: name,
		Data: data,
	})
}

// Close frees all resources allocated by the client.
// Calling close twice does nothing.
// If there is action registered for 'disconnect' message, it will be invoked.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.conn.Close()
		c.invoke("disconnect", "")
	})
	return err
}

func (c *Client) readLoop() {
	defer func() {
		if err := c.Close(); err != nil {
			c.log.Printf("error closing self: %v\n", err)
		}
	}()

	var retries = 0
	for {
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return
		}
		var msg message
		if err := c.conn.ReadJSON(&msg); err != nil {
			if retries == c.retryCount {
				return
			}
			c.log.Printf("error reading JSON from client: %v, retrying...\n", err)
			time.Sleep(c.retryInterval)
			retries++
			continue
		}
		retries = 0

		c.invoke(msg.Name, msg.Data)
	}
}

func (c *Client) invoke(name, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if action, ok := c.actions[name]; ok {
		c.log.Printf("invoking action for %s\n", name)
		action(data)
		c.log.Printf("invoking action for %s is DONE\n", name)
	}
}
