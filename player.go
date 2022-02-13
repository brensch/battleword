package battleword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Definition PlayerDefinition   `json:"definition,omitempty"`
	Games      []*PlayerGameState `json:"state,omitempty"`

	connection *PlayerConnection
}

type PlayerDefinition struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// This is all secret or not json readable types
type PlayerConnection struct {
	uri string

	// TODO: add things specific to each player
	client                      *http.Client
	concurrentConnectionLimiter chan struct{}
}

type PlayerGameState struct {
	GameID string `json:"game_id,omitempty"`

	Guesses []string `json:"guesses,omitempty"`
	Results [][]int  `json:"results,omitempty"`
	Shouts  []string `json:"shouts,omitempty"`

	Times     []time.Duration `json:"times,omitempty"`
	TotalTime time.Duration   `json:"total_time,omitempty"`
}

func InitPlayer(name, description, uri string) *Player {
	id := uuid.New()
	return &Player{
		Definition: PlayerDefinition{
			ID:          id.String(),
			Name:        name,
			Description: description,
		},
		connection: &PlayerConnection{
			uri: uri,
		},
	}
}

func (p *Player) PlayGame(g *Game) *PlayerGameState {
	gameState := &PlayerGameState{
		GameID: g.ID,
	}

	p.Games = append(p.Games, gameState)

	gameState.PlayGame(p.connection, p.Definition, g)

	return gameState

}

func (s *PlayerGameState) PlayGame(c *PlayerConnection, d PlayerDefinition, g *Game) {

	var correct bool
	var err error

	for {
		correct, err = s.DoMove(c, g.Answer)
		if err != nil {
			break
		}

		if correct {
			break
		}

		// https://i.redd.it/cw0cedsc93h81.jpg
		if len(s.Guesses) == g.numRounds {
			break
		}
	}

	for _, guessTime := range s.Times {
		s.TotalTime += guessTime
	}

	finished := "finished"
	if !correct {
		finished = "couldn't quite get it"
	}

	log.Printf("%s %s in %d turns and %d milliseconds\n", d.Name, finished, len(s.Guesses), s.TotalTime.Milliseconds())
}

func (s *PlayerGameState) DoMove(c *PlayerConnection, answer string) (bool, error) {

	start := time.Now()

	guess, err := s.GetGuess(c)
	if err != nil {
		return false, err
	}

	err = s.RecordGuess(guess, answer)
	if err != nil {
		return false, err
	}

	s.Times = append(s.Times, time.Since(start))

	correct := false
	if guess.Guess == answer {
		correct = true
	}

	return correct, nil
}

func (s *PlayerGameState) RecordGuess(guess *Guess, answer string) error {

	if !ValidGuess(guess.Guess, answer) {
		return fmt.Errorf("guess is invalid")
	}

	result := GetResult(guess.Guess, answer)

	s.Guesses = append(s.Guesses, guess.Guess)
	s.Results = append(s.Results, result)

	// TODO: also implement shouter to send shouts to everyone.
	s.Shouts = append(s.Shouts, guess.Shout)

	return nil
}

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func (s *PlayerGameState) GetGuess(c *PlayerConnection) (*Guess, error) {

	guessesJson, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/guess", c.uri), bytes.NewReader(guessesJson))
	if err != nil {
		return nil, err
	}

	// TODO: make this a proper client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var guess *Guess
	err = json.NewDecoder(res.Body).Decode(&guess)
	if err != nil {
		return nil, err
	}

	return guess, nil
}
