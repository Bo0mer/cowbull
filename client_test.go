package cowbull_test

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	. "github.com/Bo0mer/cowbull"
	"github.com/Bo0mer/cowbull/cowbullfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeNetAddr struct{}

func (fakeNetAddr) Network() string {
	return "network"
}

func (fakeNetAddr) String() string {
	return "addr"
}

var _ = Describe("Client", func() {
	var conn *cowbullfakes.FakeConn
	var c *Client

	BeforeEach(func() {
		conn = new(cowbullfakes.FakeConn)
		conn.RemoteAddrReturns(fakeNetAddr{})
	})

	AfterEach(func() {
		Expect(c.Close()).To(Succeed())
	})

	Describe("ID", func() {
		BeforeEach(func() {
			c = NewClient(conn)
		})

		It("should return a non-empty string", func() {
			Ω(c.ID()).ShouldNot(BeEmpty())
		})
	})

	Describe("SendMessage", func() {
		var err error

		JustBeforeEach(func() {
			c = NewClient(conn)
		})

		Context("when the write on connection fails", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("error writing")
				conn.WriteJSONReturns(expectedErr)
			})

			It("should have returned an error", func() {
				err = c.SendMessage("this", "will error")
				Ω(err).Should(HaveOccurred())
				Ω(err).Should(Equal(expectedErr))
			})
		})

		It("should have written to the connection", func() {
			err = c.SendMessage("name", "data")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(conn.WriteJSONCallCount()).Should(Equal(1))
		})
	})

	Describe("OnMessage", func() {
		Context("when action for a specific message is registered", func() {
			Context("and message of the same kind is recved", func() {
				var name string
				var invoked chan struct{}

				BeforeEach(func() {
					name = "Ω"
					invoked = make(chan struct{})
					var once sync.Once
					conn.ReadJSONStub = func(val interface{}) error {
						var err error
						once.Do(func() {
							err = json.Unmarshal([]byte(`{"name":"Ω"}`), val)
						})
						return err
					}
					c = NewClient(conn)
					c.OnMessage(name, func(_ string) {
						invoked <- struct{}{}
					})
				})

				It("should have been invoked", func() {
					Eventually(invoked).Should(Receive())
				})
			})
		})
	})

	Context("with broken connection", func() {
		var closed chan struct{}
		var invoked chan struct{}
		BeforeEach(func() {
			conn.ReadJSONReturns(errors.New("broken"))
			closed = make(chan struct{}, 1)
			conn.CloseStub = func() error {
				closed <- struct{}{}
				return nil
			}

			c = NewClient(conn, RetryCount(0))
			invoked = make(chan struct{}, 1)
			c.OnMessage("disconnect", func(_ string) {
				invoked <- struct{}{}
			})
		})

		It("should close it", func() {
			Eventually(closed).Should(Receive())
		})

		It("should invoke the close action", func() {
			Eventually(invoked).Should(Receive())
		})
	})

	Describe("RetryCount", func() {
		var n int
		var tries int
		var closed chan struct{}
		BeforeEach(func() {
			n = 4
			closed = make(chan struct{}, 1)
			conn.CloseStub = func() error {
				closed <- struct{}{}
				return nil
			}
			conn.ReadJSONStub = func(_ interface{}) error {
				tries++
				return errors.New("doh!")
			}
			c = NewClient(conn, RetryCount(4), RetryInterval(time.Millisecond*5))
		})

		It("should retry n times before giving up", func() {
			Eventually(closed).Should(Receive())
			// first attempt is not a retry, it is just a try.
			retries := tries - 1
			Ω(retries).Should(Equal(n))
		})
	})
})

var _ = Describe("Close", func() {

	var conn *cowbullfakes.FakeConn
	var c *Client

	BeforeEach(func() {
		conn = new(cowbullfakes.FakeConn)
		conn.RemoteAddrReturns(fakeNetAddr{})
	})
	var err error

	Context("when there is error closing the connection", func() {
		var expectedErr error
		BeforeEach(func() {
			expectedErr = errors.New("kaboom")
			conn.CloseReturns(expectedErr)
			c = NewClient(conn)

			err = conn.Close()
		})

		It("should return an error", func() {
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(expectedErr))
		})
	})

	Context("when there is action registered for disconnect", func() {
		var invoked chan struct{}

		BeforeEach(func() {
			c = NewClient(conn)
			invoked = make(chan struct{}, 1)
			c.OnMessage("disconnect", func(_ string) {
				invoked <- struct{}{}
			})

			err = c.Close()
		})

		It("should have invoked it", func() {
			Eventually(invoked).Should(Receive())
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when called twice", func() {
		BeforeEach(func() {
			c = NewClient(conn)
			Expect(c.Close()).To(Succeed())
			Expect(c.Close()).To(Succeed())
		})

		It("should close the connection only once", func() {
			Ω(conn.CloseCallCount()).Should(Equal(1))
		})
	})
})
