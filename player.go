package battleword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type PlayerState struct {
	Guesses []string `json:"guesses,omitempty"`
	Results [][]int  `json:"results,omitempty"`

	Times []time.Duration `json:"times,omitempty"`
}

type Player struct {
	Name  string      `json:"name,omitempty"`
	State PlayerState `json:"state,omitempty"`

	uri string
}

func InitPlayer(name, uri string) *Player {
	return &Player{
		Name: name,
		uri:  uri,
	}
}

func (p *Player) PlayGame(answer string, numRounds int) {

	var correct bool
	var err error

	for {
		correct, err = p.DoMove(answer)
		if err != nil {
			break
		}

		if correct {
			break
		}

		// https://i.redd.it/cw0cedsc93h81.jpg
		if len(p.State.Guesses) == numRounds {
			break
		}
	}

	var totalTime time.Duration
	for _, guessTime := range p.State.Times {
		totalTime += guessTime
	}

	finished := "finished"
	if !correct {
		finished = "couldn't quite get it"
	}

	log.Printf("%s %s in %d turns and %d milliseconds\n", p.Name, finished, len(p.State.Guesses), totalTime.Milliseconds())
}

func (p *Player) DoMove(answer string) (bool, error) {

	start := time.Now()

	guess, err := p.GetGuess()
	if err != nil {
		return false, err
	}

	err = p.State.RecordGuess(guess, answer)
	if err != nil {
		return false, err
	}

	p.State.Times = append(p.State.Times, time.Since(start))

	correct := false
	if guess == answer {
		correct = true
	}

	return correct, nil
}

func (g *PlayerState) RecordGuess(guess, answer string) error {

	if !ValidGuess(guess, answer) {
		return fmt.Errorf("guess is invalid")
	}

	result := GetResult(guess, answer)

	g.Guesses = append(g.Guesses, guess)
	g.Results = append(g.Results, result)

	return nil
}

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func (p *Player) GetGuess() (string, error) {

	guessesJson, err := json.Marshal(p.State)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/guess", p.uri), bytes.NewReader(guessesJson))
	if err != nil {
		return "", err
	}

	// TODO: make this a proper client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var guess Guess
	err = json.NewDecoder(res.Body).Decode(&guess)
	if err != nil {
		return "", err
	}

	// TODO: store this somewhere so it gets sent to other players
	log.Printf("%s shouted: %s\n", p.Name, guess.Shout)

	return guess.Guess, nil
}
