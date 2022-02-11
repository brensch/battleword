package battleword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Guesses struct {
	PreviousWords   []string `json:"previous_words,omitempty"`
	PreviousResults [][]int  `json:"previous_results,omitempty"`

	Times []time.Duration `json:"times,omitempty"`
}

type Player struct {
	Name    string  `json:"name,omitempty"`
	Guesses Guesses `json:"guesses,omitempty"`

	uri string

	// TODO: check race conditions, there's a chance we actually don't need this
	// since all calls to player will be in series
	mu sync.Mutex `json:"mu,omitempty"`
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
		if len(p.Guesses.PreviousWords) == numRounds {
			break
		}
	}

	var totalTime time.Duration
	for _, guessTime := range p.Guesses.Times {
		totalTime += guessTime
	}

	finished := "finished"
	if !correct {
		finished = "couldn't quite get it"
	}

	log.Printf("%s %s in %d turns and %d milliseconds\n", p.Name, finished, len(p.Guesses.PreviousWords), totalTime.Milliseconds())
}

func (p *Player) DoMove(answer string) (bool, error) {

	start := time.Now()

	guess, err := p.GetGuess()
	if err != nil {
		return false, err
	}

	err = p.Guesses.RecordGuess(guess, answer)
	if err != nil {
		return false, err
	}

	p.Guesses.Times = append(p.Guesses.Times, time.Since(start))

	correct := false
	if guess == answer {
		correct = true
	}

	return correct, nil
}

func (g *Guesses) RecordGuess(guess, answer string) error {

	if !ValidGuess(guess, answer) {
		return fmt.Errorf("guess is invalid")
	}

	result := GetResult(guess, answer)

	g.PreviousWords = append(g.PreviousWords, guess)
	g.PreviousResults = append(g.PreviousResults, result)

	return nil
}

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func (p *Player) GetGuess() (string, error) {

	guessesJson, err := json.Marshal(p.Guesses)
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
