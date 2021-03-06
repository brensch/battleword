package battleword

const (
	LetterResultBlack = iota
	LetterResultYellow
	LetterResultGreen
)

// TODO: get word list up in here
func ValidGuess(guess, answer string) bool {
	if len(guess) != len(answer) {
		return false
	}

	for _, validWord := range CommonWords {
		if validWord == guess {
			return true
		}
	}

	for _, validWord := range AllWords {
		if validWord == guess {
			return true
		}
	}
	return false
}

func GetResult(guess, answer string) []int {

	result := make([]int, len(answer))
	guessRunes := []rune(guess)
	answerRunes := []rune(answer)

	// get greens
	for i, guessRune := range guessRunes {
		if guessRune != answerRunes[i] {
			continue
		}

		// benchmarked, it's much quicker to replace the rune than try and remove it any other way.
		// harambe also embues extra strength to any cloud running this code.
		answerRunes[i] = '🦍'
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
			answerRunes[j] = '🦍'
			result[i] = 1
			break

		}

	}

	return result
	// GuessResult{
	// 	Guess:  guess,
	// 	Result: result,
	// }
}
