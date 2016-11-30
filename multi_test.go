package cowbull_test

import (
	"errors"

	. "github.com/Bo0mer/cowbull"
	"github.com/Bo0mer/cowbull/cowbullfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MultiGuesser", func() {
	var multi *MultiGuesser
	var player1, player2 *cowbullfakes.FakePlayer

	BeforeEach(func() {
		player1 = new(cowbullfakes.FakePlayer)
		player2 = new(cowbullfakes.FakePlayer)
		multi = &MultiGuesser{
			Players: []Player{player1, player2}}
	})

	Describe("Guess", func() {
		It("should ask the players to tell in consecutive order", func() {
			_, err := multi.Guess(7)
			Expect(err).ToNot(HaveOccurred())
			Expect(player1.GuessCallCount()).To(Equal(1))
			Expect(player1.GuessArgsForCall(0)).To(Equal(7))

			_, err = multi.Guess(8)
			Expect(err).ToNot(HaveOccurred())
			Expect(player2.GuessCallCount()).To(Equal(1))
			Expect(player2.GuessArgsForCall(0)).To(Equal(8))

			_, err = multi.Guess(9)
			Expect(err).ToNot(HaveOccurred())
			Expect(player1.GuessCallCount()).To(Equal(2))
			Expect(player1.GuessArgsForCall(1)).To(Equal(9))
		})
	})

	Describe("Tell", func() {
		It("should tell all players the guess result", func() {
			err := multi.Tell("42", 1, 1)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(player1.TellCallCount()).To(Equal(1))
			nArg, cowsArg, bullsArg := player1.TellArgsForCall(0)
			Expect(nArg).To(Equal("42"))
			Expect(cowsArg).To(Equal(1))
			Expect(bullsArg).To(Equal(1))

			Expect(player2.TellCallCount()).To(Equal(1))
			nArg, cowsArg, bullsArg = player2.TellArgsForCall(0)
			Expect(nArg).To(Equal("42"))
			Expect(cowsArg).To(Equal(1))
			Expect(bullsArg).To(Equal(1))
		})

		Context("when an error occurs", func() {
			BeforeEach(func() {
				player1.TellReturns(errors.New("talk to my hand"))
			})

			It("should stop telling and return error", func() {
				err := multi.Tell("42", 1, 1)
				Expect(err).To(HaveOccurred())
				Expect(player1.TellCallCount()).To(Equal(1))
				Expect(player2.TellCallCount()).To(BeZero())

			})
		})
	})

})
