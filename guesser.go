package cowbull

import (
	"errors"
	"strings"
)

// AIGuesser is an artificial intelligence that can guess numbers.
type AIGuesser struct {
	patterns  map[string]struct{}
	lastGuess string
}

// LocalGuesser creates new AIGuesser that guesses n-digit numbers.
func LocalGuesser(n int) *AIGuesser {
	patterns := make(map[string]struct{})
	variations([]byte("1234567890"), n, 0, patterns)
	g := &AIGuesser{
		patterns: patterns,
	}

	return g
}

// Guess returns a guess number.
func (g *AIGuesser) Guess(n int) (string, error) {
	var pattern string
	for pattern = range g.patterns {
		delete(g.patterns, pattern)
		break
	}
	g.lastGuess = pattern
	return pattern, nil
}

// Tell tells the resource of a specific guess.
func (g *AIGuesser) Tell(number string, cows, bulls int) error {
	for pattern := range g.patterns {
		c, b := computeCowsBulls(pattern, number)
		if cows != c || bulls != b {
			delete(g.patterns, pattern)
		}
	}
	if len(g.patterns) == 0 {
		return errors.New("invalid input")
	}
	return nil
}

func variations(alphabet []byte, n int, k int, out map[string]struct{}) {
	if k == n {
		out[string(alphabet[:n])] = struct{}{}
		return
	}
	for i := k; i < len(alphabet); i++ {
		alphabet[k], alphabet[i] = alphabet[i], alphabet[k]
		variations(alphabet, n, k+1, out)
		alphabet[k], alphabet[i] = alphabet[i], alphabet[k]
	}
}

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
