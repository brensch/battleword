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

func GetResult(guess, answer string) []int {

	result := make([]int, len(answer))

	// I think there's probably an optimisation i could do here if
	// i spent more time on leetcode.
	for answerPos, answerChar := range answer {
		bestPos := -1
		bestResult := 0
		for guessPos, guessChar := range guess {
			if guessChar == answerChar && guessPos == answerPos {
				bestPos = guessPos
				bestResult = 2
				break
			}

			if (bestPos == -1 || result[bestPos] > 0) && guessChar == answerChar {
				bestPos = guessPos
				bestResult = 1
			}
		}

		if bestPos >= 0 {
			result[bestPos] = bestResult
		}
	}

	return result
}
