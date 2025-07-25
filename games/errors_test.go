package games

import (
	"fmt"
	"testing"
)

// Test GameError.Error method
func TestGameError_Error(t *testing.T) {
	tests := []struct {
		input    GameError
		expected string
	}{
		{GamesIsActiveError, "game is already active"},
		{GameError(999), "unknown error"}, // Test for an unknown error code
	}

	for _, tt := range tests {
		result := tt.input.Error()
		if result != tt.expected {
			t.Errorf("GameError(%d).Error() = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

// Test WordExistsError.Error method
func TestWordExistsError_Error(t *testing.T) {
	word := "example"
	name := "testUser"
	expected := fmt.Sprintf("word '%s' already exists, added by: %s", word, name)
	err := WordExistsError{word: word, name: name}

	if result := err.Error(); result != expected {
		t.Errorf("WordExistsError.Error() = %q; want %q", result, expected)
	}
}

// Test SpecialCharError.Error method
func TestSpecialCharError_Error(t *testing.T) {
	word := "ex@mpl#"
	expected := fmt.Sprintf("word `%s` contains special chars", word)
	err := SpecialCharError{word: word}

	if result := err.Error(); result != expected {
		t.Errorf("SpecialCharError.Error() = %q; want %q", result, expected)
	}
}
