package cowbull

// MultiGuesser makes multiple guessers to look like one.
// Each guesser will be asked to guess in turn.
type MultiGuesser struct {
	Players []Player
	turn    int
}

// Guess will ask the next player for a guess and return it.
func (g *MultiGuesser) Guess(n int) (string, error) {
	if g.turn == len(g.Players) {
		g.turn = 0
	}
	number, err := g.Players[g.turn].Guess(n)
	g.turn++
	return number, err
}

// Tell tells all players the result of a guess.
// It returns on the first non-nil error.
func (g MultiGuesser) Tell(number string, cows, bulls int) error {
	for _, player := range g.Players {
		if err := player.Tell(number, cows, bulls); err != nil {
			return err
		}
	}
	return nil
}
