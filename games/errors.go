package games

import "fmt"

type GameError int

const (
	GamesIsActiveError GameError = iota
)

func (g GameError) Error() string {
	switch g {
	case GamesIsActiveError:
		return "game is already active"
	default:
		return "unknown error"
	}
}

type WordExistsError struct {
	word string
	name string
}

func (e WordExistsError) Error() string {
	return fmt.Sprintf("word '%s' already exists, added by: %s", e.word, e.name)
}

type SpecialCharError struct {
	word string
}

func (e SpecialCharError) Error() string {
	return fmt.Sprintf("word `%s` contains special chars", e.word)
}
