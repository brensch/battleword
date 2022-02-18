package battleword

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID         string             `json:"id,omitempty"`
	Definition *PlayerDefinition  `json:"definition,omitempty"`
	Games      []*PlayerGameState `json:"state,omitempty"`

	Summary *PlayerSummary `json:"player_summary,omitempty"`

	FailedToFinish bool `json:"failed_to_finish,omitempty"`

	connection *PlayerConnection
}

type PlayerSummary struct {
	TotalTime      time.Duration `json:"total_time,omitempty"`
	TotalGuesses   int           `json:"total_guesses"`
	AverageGuesses float64       `json:"average_guesses,omitempty"`
	GamesWon       int           `json:"games_won"`
	TotalVolume    float64       `json:"total_volume,omitempty"`

	Disqualified bool `json:"disqualified"`
}

type PlayerDefinition struct {
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
	Correct bool     `json:"correct"`
	Error   bool     `json:"error,omitempty"`

	shouts []string

	Times     []time.Duration `json:"times,omitempty"`
	TotalTime time.Duration   `json:"total_time,omitempty"`
}

func InitPlayer(uri string) (*Player, error) {

	client := &http.Client{
		Transport: &http.Transport{
			// odds are someone will be hosting this jankily.
			// the ramifications of a mitm attack are 0
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// need to think about setting this dynamically for humans
		Timeout: 500 * time.Millisecond,
	}

	c := &PlayerConnection{
		uri:    uri,
		client: client,
	}

	// we want the GetDefinition call to be the thing that wakes an api up if people are hosting
	// serverless. so give them a few retries just in case

	var definition *PlayerDefinition
	var err error
	for i := 0; i < 5; i++ {

		definition, err = GetDefinition(c)
		if err != nil {
			log.Printf("error getting definitions from player %s: %+v", uri, err)
			continue
		}

		break
	}

	if definition == nil {
		return nil, fmt.Errorf("failed to retrieve definition from player: %+v", err)
	}

	return &Player{
		ID:         uuid.New().String(),
		connection: c,
		Definition: definition,
	}, nil
}

func (p *Player) PlayGame(g *Game) *PlayerGameState {
	gameState := &PlayerGameState{
		GameID: g.ID,
	}

	p.Games = append(p.Games, gameState)

	gameState.PlayGame(p.connection, p.Definition, g)

	return gameState

}

func (s *PlayerGameState) PlayGame(c *PlayerConnection, d *PlayerDefinition, g *Game) {

	// var correct bool
	// var err error

	for {
		correct, err := s.DoMove(c, g.Answer)
		if err != nil {
			s.Error = true
			break
		}

		if correct {
			s.Correct = true
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

	// finished := "finished"
	// if !correct {
	// 	finished = "couldn't quite get it"
	// }

	// log.Printf("%s %s in %d turns and %d milliseconds\n", d.Name, finished, len(s.Guesses), s.TotalTime.Milliseconds())
}

func (s *PlayerGameState) DoMove(c *PlayerConnection, answer string) (bool, error) {

	start := time.Now()

	guess, err := s.GetGuess(c)
	if err != nil {
		return false, err
	}

	s.Times = append(s.Times, time.Since(start))

	err = s.RecordGuess(guess, answer)
	if err != nil {
		return false, err
	}

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
	s.shouts = append(s.shouts, guess.Shout)

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

	res, err := c.client.Do(req)
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

func (p *Player) BroadcastMatch(m *Match) error {

	matchJSON, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/results", p.connection.uri), bytes.NewReader(matchJSON))
	if err != nil {
		return err
	}

	res, err := p.connection.client.Do(req)
	if err != nil {
		return err
	}

	res.Body.Close()

	// we don't care about hearing back from the solver. it's really just us sending them the info.
	return nil

}

func GetDefinition(c *PlayerConnection) (*PlayerDefinition, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", c.uri), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var definition *PlayerDefinition
	err = json.NewDecoder(res.Body).Decode(&definition)
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (p *Player) Summarise() {

	var totalTime time.Duration
	var totalGuesses, totalGamesWon int

	for _, game := range p.Games {

		if game.Error {
			p.Summary = &PlayerSummary{
				Disqualified: true,
			}
			return
		}
		for _, guessTime := range game.Times {
			totalTime += guessTime
		}

		if game.Correct {
			totalGamesWon++
		}

		totalGuesses += len(game.Guesses)
		if !game.Correct {
			// add one if they didn't get it
			// (otherwise someone who guessed in 6 is the same as someone who failed)
			totalGuesses++
		}

	}

	p.Summary = &PlayerSummary{
		TotalTime:      totalTime,
		TotalGuesses:   totalGuesses,
		AverageGuesses: float64(totalGuesses) / float64(len(p.Games)),
		GamesWon:       totalGamesWon,
	}

}
