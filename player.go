package battleword

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Player struct {
	ID          string            `json:"id,omitempty"`
	Definition  PlayerDefinition  `json:"definition,omitempty"`
	GamesPlayed []PlayerGameState `json:"games_played,omitempty"`

	connection PlayerConnection

	mu  *sync.Mutex
	log logrus.FieldLogger
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

type GuessResult struct {
	Guess  string `json:"guess,omitempty"`
	Result []int  `json:"result,omitempty"`
}

type PlayerGameState struct {
	GameID string `json:"game_id,omitempty"`

	GuessResults []GuessResult `json:"guess_results,omitempty"`

	Correct bool   `json:"correct,omitempty"`
	Error   string `json:"error,omitempty"`

	shouts []string `json:"shouts,omitempty"`

	GuessDurationsNS []int64 `json:"guess_durations_ns,omitempty"`
}

func InitPlayer(mu *sync.Mutex, log logrus.FieldLogger, uri string) (*Player, error) {

	client := &http.Client{
		Transport: &http.Transport{
			// odds are someone will be hosting this jankily.
			// the ramifications of a mitm attack are 0
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// need to think about setting this dynamically for humans
		Timeout: 500 * time.Second,
	}

	c := PlayerConnection{
		uri:    uri,
		client: client,

		concurrentConnectionLimiter: make(chan struct{}, 100),
	}

	// we want the GetDefinition call to be the thing that wakes an api up if people are hosting
	// serverless. so give them a few retries just in case
	var definition PlayerDefinition
	var err error
	for i := 0; i < 5; i++ {
		definition, err = GetDefinition(c)
		if err != nil {
			log.
				WithFields(logrus.Fields{
					"player_uri": uri,
				}).
				WithError(err).
				Debug("failed attempt getting player definition")
			continue
		}

		break
	}

	if err != nil {
		log.
			WithFields(logrus.Fields{
				"player_uri": uri,
			}).
			WithError(err).
			Error("failed to get player definition")
		return nil, fmt.Errorf("failed to retrieve definition from player: %+v", err)
	}

	id := uuid.NewString()
	return &Player{
		ID:         id,
		connection: c,
		Definition: definition,

		mu:  mu,
		log: log.WithField("player_id", id),
	}, nil
}

func (p *Player) PlayMatch(games []Game) {

	var wgGenerate, wgListen sync.WaitGroup
	// this is buffered since we may get mutexed out of appending to the gamestate list
	// temporarily
	gameStateCHAN := make(chan PlayerGameState, 10)

	wgListen.Add(1)
	go func() {
		defer wgListen.Done()
		for game := range gameStateCHAN {
			p.mu.Lock()
			p.GamesPlayed = append(p.GamesPlayed, game)
			p.mu.Unlock()
		}
	}()

	for _, game := range games {
		wgGenerate.Add(1)
		go func(game Game) {
			defer wgGenerate.Done()
			gameStateCHAN <- PlayGame(p.connection, game)
		}(game)
	}

	wgGenerate.Wait()
	close(gameStateCHAN)
	wgListen.Wait()
}

func PlayGame(c PlayerConnection, g Game) PlayerGameState {

	state := PlayerGameState{
		GameID: g.ID,
	}
	var err error

	for {
		state, err = GetNextState(c, state, g.Answer)
		if err != nil {
			state.Error = err.Error()
			return state
		}

		if state.Correct {
			state.Correct = true
			return state
		}

		// https://i.redd.it/cw0cedsc93h81.jpg
		if len(state.GuessResults) == g.numRounds {
			return state
		}
	}
}

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func GetNextState(c PlayerConnection, s PlayerGameState, answer string) (PlayerGameState, error) {

	guessesJson, err := json.Marshal(s)
	if err != nil {
		return PlayerGameState{}, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/guess", c.uri), bytes.NewReader(guessesJson))
	if err != nil {
		return PlayerGameState{}, err
	}

	// wait for a channel to free up for this player
	c.concurrentConnectionLimiter <- struct{}{}
	defer func() { <-c.concurrentConnectionLimiter }()

	start := time.Now()
	res, err := c.client.Do(req)
	if err != nil {
		return PlayerGameState{}, err
	}
	defer res.Body.Close()

	guessDuration := time.Since(start)

	var guess Guess
	err = json.NewDecoder(res.Body).Decode(&guess)
	if err != nil {
		return PlayerGameState{}, err
	}

	if !ValidGuess(guess.Guess, answer) {
		return PlayerGameState{}, fmt.Errorf("guess is invalid: %s", guess.Guess)
	}

	result := GetResult(guess.Guess, answer)

	s.GuessResults = append(s.GuessResults, result)
	s.GuessDurationsNS = append(s.GuessDurationsNS, guessDuration.Nanoseconds())
	s.shouts = append(s.shouts, guess.Shout)

	return s, nil
}

// this struct includes the player's id to give them certainty about who they were
type PlayerMatchResults struct {
	PlayerID string `json:"player_id,omitempty"`
	Results  *Match `json:"results,omitempty"`
}

func (p *Player) BroadcastMatch(m *Match) error {

	results := PlayerMatchResults{
		PlayerID: p.ID,
		Results:  m,
	}

	matchJSON, err := json.Marshal(results)
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

func GetDefinition(c PlayerConnection) (PlayerDefinition, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", c.uri), nil)
	if err != nil {
		return PlayerDefinition{}, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return PlayerDefinition{}, err
	}

	defer res.Body.Close()

	var definition PlayerDefinition
	err = json.NewDecoder(res.Body).Decode(&definition)
	if err != nil {
		return PlayerDefinition{}, err
	}

	return definition, nil
}
