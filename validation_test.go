package battleword

import "testing"

func TestValidGuess(t *testing.T) {
	guess := "beast"
	answer := "beast"
	valid := ValidGuess(guess, answer)
	if !valid {
		t.FailNow()
	}
}
