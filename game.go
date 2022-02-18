package battleword

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID     string `json:"id,omitempty"`
	Answer string `json:"answer,omitempty"`

	Summary *GameSummary `json:"summary,omitempty"`

	numLetters int
	numRounds  int
}

type GameSummary struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`

	Fastest      Fastest      `json:"fastest,omitempty"`
	MostAccurate MostAccurate `json:"most_accurate,omitempty"`
	Loudest      Loudest      `json:"loudest,omitempty"`
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
