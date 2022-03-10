package battleword

const (
	LetterResultBlack = iota
	LetterResultYellow
	LetterResultGreen
)

// TODO: get word list up in here
func ValidGuess(guess, answer string) bool {
	return len(guess) == len(answer)
}

func GetResult(guess, answer string) GuessResult {

	result := make([]int, len(answer))
	guessRunes := []rune(guess)
	answerRunes := []rune(answer)

	// get greens
	for i, guessRune := range guessRunes {
		if guessRune != answerRunes[i] {
			continue
		}

		// benchmarked, it's much quicker to replace the rune than try and remove it any other way
		answerRunes[i] = '😒'
		result[i] = 2
	}

	// get yellows
	for i, guessRune := range guessRunes {
		if result[i] == 2 {
			continue
		}

		for j, answerRune := range answerRunes {
			if guessRune != answerRune {
				continue
			}
			answerRunes[j] = '😒'
			result[i] = 1
			break

		}

	}

	return GuessResult{
		Guess:  guess,
		Result: result,
	}
}
