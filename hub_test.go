package cowbull_test

import (
	"log"

	. "github.com/Bo0mer/cowbull"
	"github.com/Bo0mer/cowbull/cowbullfakes"
	"github.com/Bo0mer/cowbull/game"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hub", func() {
	var player *cowbullfakes.FakePlayer
	var gamer *cowbullfakes.FakeGamer
	var hub *Hub

	initHub := func() {
		logger := log.New(GinkgoWriter, "", 0)
		gamer = new(cowbullfakes.FakeGamer)
		hub = NewHub(gamer, logger)
	}

	playerWithId := func(id string) *cowbullfakes.FakePlayer {
		p := new(cowbullfakes.FakePlayer)
		p.IDReturns(id)
		return p
	}

	Describe("Add", func() {
		var announced chan struct{}
		BeforeEach(func() {
			announced = make(chan struct{}, 1)

			player = playerWithId("random")
			player.AnnouncePlayersStub = func(ps []PlayerEntry) error {
				Expect(len(ps)).To(Equal(1))
				Expect(ps[0].ID).To(Equal("random"))

				announced <- struct{}{}

				return nil
			}

			initHub()
			hub.Add(player)
		})

		It("should announce all players to the added one", func() {
			Eventually(announced).Should(Receive())
		})
	})

	Describe("Remove", func() {
		var announced chan struct{}
		BeforeEach(func() {
			announced = make(chan struct{}, 1)

			player1 := playerWithId("random1")
			player2 := playerWithId("random2")

			initHub()
			hub.Add(player1)

			call := 0
			player2.AnnouncePlayersStub = func(ps []PlayerEntry) error {
				// We want to check the second announcement, since the first
				// one will be when the player joins.
				if call == 1 {
					Expect(len(ps)).To(Equal(1))
					Expect(ps[0].ID).To(Equal("random2"))
				}
				call++
				announced <- struct{}{}

				return nil
			}
			hub.Add(player2)
			// Ensure that the announcement for joining is received.
			Eventually(announced).Should(Receive())
			hub.Remove("random1")
		})

		It("should announce that a player left", func() {
			Eventually(announced).Should(Receive())
		})
	})

	Describe("NewGame", func() {
		var settings GameSettings
		var g *game.Game
		var err error
		BeforeEach(func() {
			player = new(cowbullfakes.FakePlayer)
			initHub()
		})

		JustBeforeEach(func() {
			g, err = hub.NewGame(player, settings)
		})

		Context("with an invalid settings", func() {
			BeforeEach(func() {
				settings.Role = "hacker"
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(g).To(BeNil())
			})
		})

		Context("with AI thinker and player guesser", func() {
			BeforeEach(func() {
				settings = GameSettings{
					Role: RoleGuesser,
					AI:   true,
				}
			})

			It("should have created specific game", func() {
				Expect(gamer.GameCallCount()).To(Equal(1))
				t, g := gamer.GameArgsForCall(0)
				_, ok := t.(*AIThinker)
				Expect(ok).To(BeTrue())

				_, ok = g.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
			})
		})

		Context("with request for guesser for two players", func() {
			// TODO(ivan): DRY this out
			BeforeEach(func() {
				settings = GameSettings{
					Role:      RoleGuesser,
					Opponents: []string{"random2"},
				}
				player = playerWithId("random1")
				player2 := playerWithId("random2")

				initHub()
				hub.Add(player)
				hub.Add(player2)
			})

			It("should have created specific game", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gamer.GameCallCount()).To(Equal(1))

				t, g := gamer.GameArgsForCall(0)
				fp, ok := t.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
				Expect(fp.ID()).To(Equal("random2"))

				fp, ok = g.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
				Expect(fp.ID()).To(Equal("random1"))
			})
		})

		Context("with request for thinker for two players", func() {
			BeforeEach(func() {
				settings = GameSettings{
					Role:      RoleThinker,
					Opponents: []string{"random2"},
				}
				player = playerWithId("random1")
				player2 := playerWithId("random2")

				initHub()
				hub.Add(player)
				hub.Add(player2)
			})

			It("should have created specific game", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gamer.GameCallCount()).To(Equal(1))

				t, g := gamer.GameArgsForCall(0)
				fp, ok := t.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
				Expect(fp.ID()).To(Equal("random1"))

				fp, ok = g.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
				Expect(fp.ID()).To(Equal("random2"))
			})
		})

		Context("with request for thinker for >2 players", func() {
			BeforeEach(func() {
				settings = GameSettings{
					Role:      RoleThinker,
					Opponents: []string{"random2", "random3"},
				}
				player = playerWithId("random1")
				player2 := playerWithId("random2")
				player3 := playerWithId("random3")

				initHub()
				hub.Add(player)
				hub.Add(player2)
				hub.Add(player3)
			})

			It("should have created game with multiguesser", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gamer.GameCallCount()).To(Equal(1))

				t, g := gamer.GameArgsForCall(0)
				fp, ok := t.(*cowbullfakes.FakePlayer)
				Expect(ok).To(BeTrue())
				Expect(fp.ID()).To(Equal("random1"))

				_, ok = g.(*MultiGuesser)
				Expect(ok).To(BeTrue())
			})

		})
	})
})
