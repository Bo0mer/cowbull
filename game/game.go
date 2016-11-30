package game

//go:generate counterfeiter . Thinker
//go:generate counterfeiter . Guesser

// Thinker represents a player that thinks of a number and answers questions
// about it.
type Thinker interface {
	// Number makes the player to think of a number and returns the count of
	// its digits.
	Think() (int, error)
	// Try should return the cows and bulls in the specified number.
	// If the number's digit count is not equal to the one returned by Number,
	// the player should return (0, 0).
	Try(string) (cows int, bulls int, err error)
}

// Guesser represents a player that tries to guess a number.
type Guesser interface {
	// Guess should return a guess number consisting of n digits.
	Guess(n int) (string, error)
	// Tell should tell the player the result of his guess.
	Tell(string, int, int) error
}

// Game represents a cowbull game.
type Game struct {
	thinker Thinker
	guesser Guesser
}

// NewGame creates new gime with the provided players.
func New(thinker Thinker, guesser Guesser) *Game {
	return &Game{
		thinker: thinker,
		guesser: guesser,
	}
}

// Play plays the game with the players.
func (g *Game) Play() error {
	var err error
	var digits int
	if digits, err = g.thinker.Think(); err != nil {
		return err
	}
	for {
		guess, err := g.guesser.Guess(digits)
		if err != nil {
			return err
		}

		cows, bulls, err := g.thinker.Try(guess)
		if err != nil {
			return err
		}

		if err = g.guesser.Tell(guess, cows, bulls); err != nil {
			return err
		}

		if digits == bulls {
			return nil
		}
	}
}
