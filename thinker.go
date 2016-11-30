package cowbull

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type AIThinker struct {
	digits int
	number string

	perm func(int) []int
}

func LocalThinker(digits int) *AIThinker {
	return &AIThinker{
		digits: digits,
		perm:   rand.Perm,
	}
}

func NewLocalThinker(digits int, perm func(int) []int) *AIThinker {
	return &AIThinker{
		digits: digits,
		perm:   perm,
	}
}

func (p *AIThinker) Think() (int, error) {
	var err error
	p.number, err = p.generateNumber()
	if err != nil {
		return 0, err
	}
	return len(p.number), nil
}

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
