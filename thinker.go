package cowbull

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// AIThinker is an artificial intelligence that can think of numbers.
type AIThinker struct {
	digits int
	number string

	perm func(int) []int
}

// LocalThinker creates new AIThinker that thinks only of n-digit numbers.
func LocalThinker(n int) *AIThinker {
	return &AIThinker{
		digits: n,
		perm:   rand.Perm,
	}
}

// NewLocalThinker creates new AIThinker that thinks only of n-digit numbers.
// The perm argument specifies a function that generates permutations in the
// range [0, n). For more info, see the doc of rand.Perm.
// This method is used for testing. If you're using it, you're doing something
// wrong or you're trying to be clever. Remember: clear is better than clever.
func NewLocalThinker(n int, perm func(int) []int) *AIThinker {
	return &AIThinker{
		digits: n,
		perm:   perm,
	}
}

// Think thinks of a number.
func (p *AIThinker) Think() (int, error) {
	var err error
	p.number, err = p.generateNumber()
	if err != nil {
		return 0, err
	}
	return len(p.number), nil
}

// Try returns the cows and bulls for number.
func (p *AIThinker) Try(number string) (int, int, error) {
	if len(number) != len(p.number) {
		return 0, 0, errors.New("local player: try number digit count mismatch")
	}
	cows, bulls := 0, 0
	for i := range number {
		if number[i] == p.number[i] {
			bulls++
			continue
		}
		if strings.Contains(p.number, string(number[i])) {
			cows++
		}
	}
	return cows, bulls, nil
}

func (p *AIThinker) generateNumber() (string, error) {
	if p.digits < 1 || p.digits > 10 {
		return "", fmt.Errorf("invalid digit count %d", p.digits)
	}
	perm := p.perm(10)
	if perm[0] == 0 {
		randIdx := rand.Intn(p.digits-1) + 1
		perm[0], perm[randIdx] = perm[randIdx], perm[0]
	}
	res := ""
	for i := 0; i < p.digits; i++ {
		res = fmt.Sprintf("%s%d", res, perm[i])
	}
	return res, nil
}
