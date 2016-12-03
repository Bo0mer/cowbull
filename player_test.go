package cowbull_test

import (
	"fmt"
	"time"

	. "github.com/Bo0mer/cowbull"
	"github.com/Bo0mer/cowbull/cowbullfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemotePlayer", func() {

	var messenger *cowbullfakes.FakeMessenger
	var player *RemotePlayer

	BeforeEach(func() {
		messenger = new(cowbullfakes.FakeMessenger)
	})

	Describe("NewRemotePlayer", func() {
		var subscriptions map[string]struct{}

		itShouldSubscribeFor := func(kind string) {
			It(fmt.Sprintf("should subscribe for %s", kind), func() {
				_, ok := subscriptions[kind]
				Expect(ok).To(BeTrue())
			})
		}

		BeforeEach(func() {
			subscriptions = make(map[string]struct{})
			messenger.OnMessageStub = func(kind string, _ func(string)) {
				subscriptions[kind] = struct{}{}
			}
			player = NewRemotePlayer(messenger, time.Second)
		})

		itShouldSubscribeFor("name")

		itShouldSubscribeFor("think")

		itShouldSubscribeFor("guess")

		itShouldSubscribeFor("try")
	})

	Describe("ID", func() {
		BeforeEach(func() {
			messenger.IDReturns("id")
			player = NewRemotePlayer(messenger, time.Second)
		})

		It("should return the messenger's ID", func() {
			Expect(player.ID()).To(Equal("id"))
		})
	})

	Describe("AnnouncePlayers", func() {
		var err error
		BeforeEach(func() {
			player = NewRemotePlayer(messenger, time.Second)
			err = player.AnnouncePlayers([]PlayerEntry{{ID: "1", Name: "Bob"}})
		})

		It("should send a 'players' message", func() {
			Expect(messenger.SendMessageCallCount()).To(Equal(1))
			argKind, argData := messenger.SendMessageArgsForCall(0)
			Expect(argKind).To(Equal("players"))
			Expect(argData).To(Equal(`[{"id":"1","name":"Bob"}]`))
		})

		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Think", func() {
		var expectedDigits int
		var digits int
		var err error

		JustBeforeEach(func() {
			player = NewRemotePlayer(messenger, time.Second)
			digits, err = player.Think()
		})

		Context("when the think response arrives on time", func() {
			BeforeEach(func() {
				expectedDigits = 4
				messenger.OnMessageStub = func(kind string, action func(data string)) {
					if kind == "think" {
						action(fmt.Sprintf(`{"digits":%d}`, expectedDigits))
					}
				}
			})

			It("should send a 'think' message", func() {
				Expect(messenger.SendMessageCallCount()).To(Equal(1))
				argKind, argData := messenger.SendMessageArgsForCall(0)
				Expect(argKind).To(Equal("think"))
				Expect(argData).To(BeEmpty())
			})

			It("should return the number of digits given by the messenger", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(digits).To(Equal(expectedDigits))
			})
		})

		Context("when there is no think response", func() {
			It("should return a timed out error", func() {
				Ω(digits).Should(BeZero())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(Equal("remoteplayer: think timed out"))
			})
		})
	})

	Describe("Guess", func() {
		var digits int
		var expectedNumber string
		var number string
		var err error

		JustBeforeEach(func() {
			player = NewRemotePlayer(messenger, time.Millisecond*10)
			number, err = player.Guess(digits)
		})

		Context("when the guess response arrives on time", func() {
			BeforeEach(func() {
				expectedNumber = "4201"
				digits = len(expectedNumber)
				messenger.OnMessageStub = func(kind string, action func(data string)) {
					if kind == "guess" {
						action(fmt.Sprintf(`{"number":"%s"}`, expectedNumber))
					}
				}
			})

			It("should return the number given by the messenger", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(number).To(Equal(expectedNumber))
			})

			It("should send a 'guess' message", func() {
				Expect(messenger.SendMessageCallCount()).To(Equal(1))
				argKind, argData := messenger.SendMessageArgsForCall(0)
				Expect(argKind).To(Equal("guess"))
				Expect(argData).To(Equal(fmt.Sprintf(`{"digits":%d}`, digits)))
			})
		})

		Context("when there is no guess response", func() {
			It("should return a timed out error", func() {
				Ω(number).Should(BeEmpty())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(Equal("remoteplayer: guess timed out"))
			})
		})
	})

	Describe("Try", func() {
		var guess string
		var cows, bulls int
		var expectedCows, expectedBulls int
		var err error

		JustBeforeEach(func() {
			player = NewRemotePlayer(messenger, time.Second)
			cows, bulls, err = player.Try(guess)
		})

		Context("when the try response arrives on time", func() {
			BeforeEach(func() {
				guess = "4201"
				expectedCows = 2
				expectedBulls = 2
				messenger.OnMessageStub = func(kind string, action func(data string)) {
					if kind == "try" {
						action(fmt.Sprintf(`{"cows":%d,"bulls":%d}`, expectedCows, expectedBulls))
					}
				}
			})
			It("should send a 'try' message", func() {
				Expect(messenger.SendMessageCallCount()).To(Equal(1))
				argKind, argData := messenger.SendMessageArgsForCall(0)
				Expect(argKind).To(Equal("try"))
				Expect(argData).To(Equal(fmt.Sprintf(`{"number":"%s"}`, guess)))
			})

			It("should return the cows & bulls given by the messenger", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(cows).To(Equal(expectedCows))
				Expect(bulls).To(Equal(expectedBulls))
			})
		})

		Context("when there is no try response", func() {
			It("should return a timed out error", func() {
				Ω(cows).Should(BeZero())
				Ω(bulls).Should(BeZero())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(Equal("remoteplayer: try timed out"))
			})
		})

	})

	Describe("Tell", func() {
		var number string
		var cows, bulls int
		var err error
		BeforeEach(func() {
			number = "4201"
			cows = 2
			bulls = 2
			player = NewRemotePlayer(messenger, time.Second)
			err = player.Tell(number, cows, bulls)
		})

		It("should send a 'tell' message", func() {
			Expect(messenger.SendMessageCallCount()).To(Equal(1))
			argKind, argData := messenger.SendMessageArgsForCall(0)
			Expect(argKind).To(Equal("tell"))
			expected := fmt.Sprintf(`{"number":"%s","cows":%d,"bulls":%d}`, number, cows, bulls)
			Expect(argData).To(Equal(expected))
		})

		It("should not return an error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Name", func() {
		Context("when name is not set", func() {
			BeforeEach(func() {
				player = NewRemotePlayer(messenger, time.Second)
			})
			It("should return empty string", func() {
				Expect(player.Name()).To(BeEmpty())
			})
		})

		Context("when set by the messenger", func() {
			var name string
			var expectedName string
			BeforeEach(func() {
				expectedName = "Alice"
				messenger.OnMessageStub = func(kind string, action func(data string)) {
					if kind == "name" {
						action(fmt.Sprintf(`{"name":"%s"}`, expectedName))
					}
				}
				player = NewRemotePlayer(messenger, time.Second)
				name = player.Name()
			})

			It("should return the value provided by messenger", func() {
				Expect(name).To(Equal(expectedName))
			})
		})
	})

})
