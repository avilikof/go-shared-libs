package games

import (
	"log/slog"
	"strings"
	"unicode"
)

type FirstLetterGame struct {
	storage    map[string]string
	isActive   bool
	players    []string
	nameOffset int8
}

func StartFirstLetterGame(players []string) *FirstLetterGame {
	return &FirstLetterGame{
		isActive:   false,
		storage:    make(map[string]string),
		players:    players,
		nameOffset: 0,
	}
}

func (f *FirstLetterGame) IsActive() bool {
	return f.isActive
}

func (f *FirstLetterGame) WordCount() int {
	return len(f.storage)
}
func (f *FirstLetterGame) Write(word, name string) error {

	// Word should not contain special chars
	if containsSpecialChars(word) {
		err := SpecialCharError{word: word}
		slog.Error(err.Error())
		return err
	}

	// prior to writing word must be formatted
	word = applyFormat(word)

	// If word already stored name of the person must be returned
	// Else add word and name to the storage
	_, storedName, exists := f.get(word)
	if exists {
		err := WordExistsError{word: word, name: storedName}
		slog.Error(err.Error())
		return err
	}
	f.put(word, name)
	return nil
}

// Assign `storage` with empty map and set `isActive` to `false`
func (f *FirstLetterGame) Clear() error {
	f.storage = make(map[string]string)
	f.isActive = false
	return nil
}
func (f *FirstLetterGame) GetPlayers() []string {
	return f.players
}
func (f *FirstLetterGame) GetWords() *[]string {
	var words []string
	for key := range f.storage {
		words = append(words, key)
	}
	return &words
}
func (f *FirstLetterGame) NextName() string {
	if f.nameOffset < int8(len(f.players)) {
		f.nameOffset += 1
	}
	if f.nameOffset == int8(len(f.players)) {
		f.nameOffset = 0
	}
	return f.CurrentName()
}
func (f *FirstLetterGame) CurrentName() string {
	return f.players[f.nameOffset]
}

func (f *FirstLetterGame) get(word string) (string, string, bool) {
	name, exists := f.storage[word]
	if exists {
		return word, name, exists
	}
	return "", "", false
}

func (f *FirstLetterGame) put(word, name string) {
	f.storage[word] = name
}

func containsSpecialChars(word string) bool {
	for _, letter := range word {
		if !unicode.IsLetter(letter) {
			return true
		}
	}
	return false
}

func applyFormat(word string) string {
	word = strings.TrimSpace(word)
	word = strings.ToLower(word)
	return word
}
