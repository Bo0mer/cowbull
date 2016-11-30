package cowbull_test

import (
	"fmt"

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
			player = NewRemotePlayer(messenger)
		})

		itShouldSubscribeFor("name")

		itShouldSubscribeFor("think")

		itShouldSubscribeFor("guess")

		itShouldSubscribeFor("try")
	})

	Describe("ID", func() {
		BeforeEach(func() {
			messenger.IDReturns("id")
			player = NewRemotePlayer(messenger)
		})

		It("should return the messenger's ID", func() {
			Expect(player.ID()).To(Equal("id"))
		})
	})

	Describe("AnnouncePlayers", func() {
		BeforeEach(func() {
			player = NewRemotePlayer(messenger)
			player.AnnouncePlayers([]PlayerEntry{{ID: "1", Name: "Bob"}})
		})

		It("should send a 'players' message", func() {
			Expect(messenger.SendMessageCallCount()).To(Equal(1))
			argKind, argData := messenger.SendMessageArgsForCall(0)
			Expect(argKind).To(Equal("players"))
			Expect(argData).To(Equal(`[{"id":"1","name":"Bob"}]`))
		})
	})

	Describe("Think", func() {
		var expectedDigits int
		var digits int
		var err error
		BeforeEach(func() {
			expectedDigits = 4
			messenger.OnMessageStub = func(kind string, action func(data string)) {
				if kind == "think" {
					action(fmt.Sprintf(`{"digits":%d}`, expectedDigits))
				}
			}
			player = NewRemotePlayer(messenger)
			digits, err = player.Think()
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

	Describe("Guess", func() {
		var digits int
		var expectedNumber string
		var number string
		var err error
		BeforeEach(func() {
			expectedNumber = "4201"
			digits = len(expectedNumber)
			messenger.OnMessageStub = func(kind string, action func(data string)) {
				if kind == "guess" {
					action(fmt.Sprintf(`{"number":"%s"}`, expectedNumber))
				}
			}
			player = NewRemotePlayer(messenger)
			number, err = player.Guess(digits)
		})

		It("should send a 'guess' message", func() {
			Expect(messenger.SendMessageCallCount()).To(Equal(1))
			argKind, argData := messenger.SendMessageArgsForCall(0)
			Expect(argKind).To(Equal("guess"))
			Expect(argData).To(Equal(fmt.Sprintf(`{"digits":%d}`, digits)))
		})

		It("should return the number given by the messenger", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(number).To(Equal(expectedNumber))
		})
	})

	Describe("Try", func() {
		var guess string
		var cows, bulls int
		var expectedCows, expectedBulls int
		var err error
		BeforeEach(func() {
			guess = "4201"
			expectedCows = 2
			expectedBulls = 2
			messenger.OnMessageStub = func(kind string, action func(data string)) {
				if kind == "try" {
					action(fmt.Sprintf(`{"cows":%d,"bulls":%d}`, expectedCows, expectedBulls))
				}
			}
			player = NewRemotePlayer(messenger)
			cows, bulls, err = player.Try(guess)
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

	Describe("Tell", func() {
		var number string
		var cows, bulls int
		BeforeEach(func() {
			number = "4201"
			cows = 2
			bulls = 2
			player = NewRemotePlayer(messenger)
			player.Tell(number, cows, bulls)
		})

		It("should send a 'tell' message", func() {
			Expect(messenger.SendMessageCallCount()).To(Equal(1))
			argKind, argData := messenger.SendMessageArgsForCall(0)
			Expect(argKind).To(Equal("tell"))
			expected := fmt.Sprintf(`{"number":"%s","cows":%d,"bulls":%d}`, number, cows, bulls)
			Expect(argData).To(Equal(expected))
		})
	})

	Describe("Name", func() {
		Context("when name is not set", func() {
			BeforeEach(func() {
				player = NewRemotePlayer(messenger)
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
				player = NewRemotePlayer(messenger)
				name = player.Name()
			})

			It("should return the value provided by messenger", func() {
				Expect(name).To(Equal(expectedName))
			})
		})
	})

})
