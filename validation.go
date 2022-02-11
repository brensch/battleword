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

	// TODO: make only the correct number of yellows show up
	for guessCharacterPosition, guessCharacter := range guess {
		result[guessCharacterPosition] = LetterResultBlack
		for answerCharacterPosition, answerCharacter := range answer {
			if guessCharacter == answerCharacter && guessCharacterPosition == answerCharacterPosition {
				result[guessCharacterPosition] = LetterResultGreen
				break
			}

			if guessCharacter == answerCharacter {
				result[guessCharacterPosition] = LetterResultYellow
			}

		}
	}

	return result
}
