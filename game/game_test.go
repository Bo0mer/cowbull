package game_test

import (
	"errors"

	. "github.com/Bo0mer/cowbull/game"
	"github.com/Bo0mer/cowbull/game/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {

	var thinker *fakes.FakeThinker
	var guesser *fakes.FakeGuesser
	var game *Game

	var err error
	var expectedErr error

	itShouldHaveErrored := func() {
		It("should have errored", func() {
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(expectedErr))
		})
	}

	BeforeEach(func() {
		thinker = new(fakes.FakeThinker)
		guesser = new(fakes.FakeGuesser)
	})

	JustBeforeEach(func() {
		game = New(thinker, guesser)
		err = game.Play()
	})

	Describe("Play", func() {
		Context("when the thinker fails on thinking", func() {
			BeforeEach(func() {
				expectedErr = errors.New("error thinking")
				thinker.ThinkReturns(0, expectedErr)
			})

			itShouldHaveErrored()

		})

		Context("when the thinker thinks of a number", func() {
			var digits int
			BeforeEach(func() {
				digits = 2
				thinker.ThinkReturns(digits, nil)
				thinker.TryReturns(0, 2, nil)
			})

			It("should ask the guesser to guess it", func() {
				Ω(guesser.GuessCallCount()).Should(Equal(1))
				Ω(guesser.GuessArgsForCall(0)).Should(Equal(digits))
			})

			Context("when the guesser errors on guessing", func() {
				BeforeEach(func() {
					expectedErr = errors.New("error guessing")
					guesser.GuessReturns("", expectedErr)
				})

				itShouldHaveErrored()
			})

			Context("when the guesser returns a guess", func() {
				var guess string
				BeforeEach(func() {
					guess = "42"
					guesser.GuessReturns(guess, nil)
				})

				It("should try the number with the thinker", func() {
					Ω(thinker.TryCallCount()).Should(Equal(1))
					Ω(thinker.TryArgsForCall(0)).Should(Equal(guess))
				})

				Context("when guess trial fails", func() {
					BeforeEach(func() {
						expectedErr = errors.New("error try")
						thinker.TryReturns(0, 0, expectedErr)
					})

					itShouldHaveErrored()
				})

				Context("when guess trial succeeds", func() {
					var cows, bulls int
					BeforeEach(func() {
						cows = 0
						bulls = 2
						thinker.TryReturns(cows, bulls, nil)
					})

					It("should tell the result to the guesser", func() {
						Ω(guesser.TellCallCount()).Should(Equal(1))
						argCows, argBulls := guesser.TellArgsForCall(0)
						Ω(argCows).Should(Equal(cows))
						Ω(argBulls).Should(Equal(bulls))
					})

					It("should end the game if the number is guessed", func() {
						Ω(err).ShouldNot(HaveOccurred())
					})

					Context("when telling the guesser fails", func() {
						BeforeEach(func() {
							expectedErr = errors.New("tell me baby one more time")
							guesser.TellReturns(expectedErr)
						})

						itShouldHaveErrored()
					})
				})
			})
		})
	})
})
