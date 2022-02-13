package battleword

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID     string `json:"id,omitempty"`
	Answer string `json:"answer,omitempty"`

	Result *GameResult `json:"result,omitempty"`

	numLetters int
	numRounds  int
}

type GameResult struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`

	Fastest      FastestPlayer      `json:"fastest,omitempty"`
	MostAccurate MostAccuratePlayer `json:"most_accurate,omitempty"`
	Loudest      LoudestPlayer      `json:"loudest,omitempty"`
}

type FastestPlayer struct {
	Player PlayerDefinition `json:"player,omitempty"`
	Time   time.Duration    `json:"time,omitempty"`
}

type MostAccuratePlayer struct {
	Player             PlayerDefinition `json:"player,omitempty"`
	AverageGuessLength float64          `json:"average_guess_length,omitempty"`
}

type LoudestPlayer struct {
	Player PlayerDefinition `json:"player,omitempty"`
	Volume float64          `json:"volume,omitempty"`
}

func InitGame(commonWords []string, numLetters, numRounds int) *Game {

	id := uuid.New()
	game := &Game{
		ID: id.String(),

		// TODO:make this adjust to numletters
		Answer: GetRandomWord(commonWords),

		numLetters: numLetters,
		numRounds:  numRounds,
	}

	return game
}

func GetRandomWord(words []string) string {
	// don't need to do this any more legit than this i don't think
	rand.Seed(time.Now().UnixNano())

	return words[rand.Intn(len(words))]
}
