package games

import (
	"testing"
)

// Test StartFirstLetterGame
func TestStartFirstLetterGame(t *testing.T) {
	players := []string{"Alice", "Bob"}
	game := StartFirstLetterGame(players)

	if game.IsActive() {
		t.Errorf("Expected game to be inactive initially")
	}
	if len(game.storage) != 0 {
		t.Errorf("Expected empty storage, got %v", game.storage)
	}
	if len(game.GetPlayers()) != len(players) {
		t.Errorf("Expected %d players, got %d", len(players), len(game.GetPlayers()))
	}
}

// Test IsActive method
func TestIsActive(t *testing.T) {
	game := StartFirstLetterGame([]string{"Alice", "Bob"})
	if game.IsActive() {
		t.Errorf("Expected IsActive to be false initially")
	}
}

// Test WordCount method
func TestWordCount(t *testing.T) {
	game := StartFirstLetterGame([]string{"Alice", "Bob"})
	game.storage["test"] = "Alice"

	if count := game.WordCount(); count != 1 {
		t.Errorf("Expected WordCount to be 1, got %d", count)
	}
}

// Test Write method
func TestWrite(t *testing.T) {
	game := StartFirstLetterGame([]string{"Alice", "Bob"})

	// Test with special characters
	err := game.Write("hello!", "Alice")
	if err == nil || err.Error() != "word `hello!` contains special chars" {
		t.Errorf("Expected SpecialCharError for 'hello!', got %v", err)
	}

	// Test writing a word for the first time
	err = game.Write("hello", "Alice")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test writing the same word by another player
	err = game.Write("hello", "Bob")
	if err == nil || err.Error() != "word 'hello' already exists, added by: Alice" {
		t.Errorf("Expected WordExistsError for 'hello', got %v", err)
	}
}

// Test Clear method
func TestClear(t *testing.T) {
	game := StartFirstLetterGame([]string{"Alice", "Bob"})
	game.Write("hello", "Alice")
	game.Clear()

	if game.WordCount() != 0 {
		t.Errorf("Expected WordCount to be 0 after Clear, got %d", game.WordCount())
	}
	if game.IsActive() {
		t.Errorf("Expected game to be inactive after Clear")
	}
}

// Test GetPlayers method
func TestGetPlayers(t *testing.T) {
	players := []string{"Alice", "Bob"}
	game := StartFirstLetterGame(players)

	if len(game.GetPlayers()) != len(players) {
		t.Errorf("Expected %d players, got %d", len(players), len(game.GetPlayers()))
	}
}

// Test GetWords method
func TestGetWords(t *testing.T) {
	game := StartFirstLetterGame([]string{"Alice", "Bob"})
	game.Write("hello", "Alice")
	game.Write("world", "Bob")

	words := game.GetWords()
	if len(*words) != 2 {
		t.Errorf("Expected 2 words, got %d", len(*words))
	}
	if (*words)[0] != "hello" && (*words)[1] != "world" {
		t.Errorf("Expected words to contain 'hello' and 'world', got %v", *words)
	}
}

// Test NextName and CurrentName methods
func TestNextNameAndCurrentName(t *testing.T) {
	players := []string{"Alice", "Bob", "Charlie"}
	game := StartFirstLetterGame(players)

	// Test initial CurrentName
	if game.CurrentName() != "Alice" {
		t.Errorf("Expected CurrentName to be 'Alice', got '%s'", game.CurrentName())
	}

	// Test NextName progression
	if game.NextName() != "Bob" {
		t.Errorf("Expected NextName to be 'Bob', got '%s'", game.CurrentName())
	}
	if game.NextName() != "Charlie" {
		t.Errorf("Expected NextName to be 'Charlie', got '%s'", game.CurrentName())
	}

	// Test wrap-around to the first player
	if game.NextName() != "Alice" {
		t.Errorf("Expected NextName to wrap around to 'Alice', got '%s'", game.CurrentName())
	}
}

// Test containsSpecialChars function
func TestContainsSpecialChars(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"hello", false},
		{"hello!", true},
		{"123", true},
		{"hello_world", true},
	}

	for _, tt := range tests {
		if result := containsSpecialChars(tt.word); result != tt.expected {
			t.Errorf("containsSpecialChars(%s) = %v; want %v", tt.word, result, tt.expected)
		}
	}
}

// Test applyFormat function
func TestApplyFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" Hello ", "hello"},
		{"HELLO", "hello"},
		{"  HeLLo World  ", "hello world"},
	}

	for _, tt := range tests {
		if result := applyFormat(tt.input); result != tt.expected {
			t.Errorf("applyFormat(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
