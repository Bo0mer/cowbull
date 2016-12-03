package cowbull_test

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	. "github.com/Bo0mer/cowbull"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AIGuesser", func() {
	var number string
	var g *AIGuesser
	var seed int64

	BeforeEach(func() {
		seed = time.Now().UnixNano()
		rand.Seed(seed)
		// make test reproducable
		fmt.Fprintf(GinkgoWriter, "AIGuesser: using seed %d\n", seed)
	})

	Context("when the number has four digits", func() {
		BeforeEach(func() {
			g = LocalGuesser(4)
			number = "4201"
		})

		It("should be guessed within 8 tries", func() {
			i := 0
			for {
				i++
				guess, err := g.Guess(4)
				Ω(err).ShouldNot(HaveOccurred())
				if guess == number {
					break
				}
				c, b := computeCowsBulls(number, guess)
				g.Tell(guess, c, b)
			}
			Ω(i).Should(BeNumerically("<", 8))
		})
	})

	Context("when the thinker is speculating", func() {
		BeforeEach(func() {
			g = LocalGuesser(4)
			number = "4201"
		})

		It("should return an error", func() {
			var err error
			for err == nil {
				var guess string
				guess, err = g.Guess(4)
				Ω(err).ShouldNot(HaveOccurred())
				err = g.Tell(guess, rand.Intn(4), rand.Intn(4))
			}
			Ω(err).Should(HaveOccurred())
		})
	})
})

func computeCowsBulls(origin, guess string) (int, int) {
	var cows, bulls int
	for i := range guess {
		if guess[i] == origin[i] {
			bulls++
			continue
		}
		if strings.Contains(origin, string(guess[i])) {
			cows++
		}
	}
	return cows, bulls
}
