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
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	ConcurrentConnLimit int    `json:"concurrent_connection_limit,omitempty"`
}

// This is all secret or not json readable types
type PlayerConnection struct {
	uri string

	// TODO: add things specific to each player
	client                      *http.Client
	concurrentConnectionLimiter chan struct{}
}

const (
	GuessIDHeader = "guessID"
)

type GuessResult struct {
	ID string `json:"id,omitempty"`

	Guess  string `json:"guess,omitempty"`
	Result []int  `json:"result,omitempty"`
}

type PlayerGameState struct {
	GameID string `json:"game_id,omitempty"`

	GuessResults []GuessResult `json:"guess_results,omitempty"`

	Correct bool   `json:"correct,omitempty"`
	Error   string `json:"error,omitempty"`

	shouts []string

	GuessDurationsNS []int64 `json:"guess_durations_ns,omitempty"`
}

func InitPlayer(mu *sync.Mutex, log logrus.FieldLogger, uri string) (*Player, error) {

	// we want the GetDefinition call to be the thing that wakes an api up if people are hosting
	// serverless. so give them a few retries just in case
	var definition PlayerDefinition
	var err error
	for i := 0; i < 5; i++ {
		definition, err = GetDefinition(uri)
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

	client := &http.Client{
		Transport: &http.Transport{
			// odds are someone will be hosting this jankily.
			// the ramifications of a mitm attack are 0
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// Need to think about setting this dynamically for humans.
		// We are getting the definition of the player before this so could easily figure out somehow from there.
		Timeout: 500 * time.Second,
	}
	c := PlayerConnection{
		uri:    uri,
		client: client,
	}

	// Set the connection limit based off what the player specified
	connectionLimit := definition.ConcurrentConnLimit
	if connectionLimit == 0 {
		connectionLimit = 5
	}
	c.concurrentConnectionLimiter = make(chan struct{}, connectionLimit)

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

	for {
		state = GetNextState(c, state, g.Answer)

		// https://i.redd.it/cw0cedsc93h81.jpg
		if state.Correct || state.Error != "" || len(state.GuessResults) == g.numRounds {
			return state
		}
	}
}

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func GetNextState(c PlayerConnection, s PlayerGameState, answer string) PlayerGameState {

	guessesJson, err := json.Marshal(s)
	if err != nil {
		s.Error = err.Error()
		return s
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/guess", c.uri), bytes.NewReader(guessesJson))
	if err != nil {
		s.Error = err.Error()
		return s
	}

	id := uuid.New().String()

	req.Header.Add(GuessIDHeader, id)

	// Make sure we don't go over the concurrent connection limit for this player.
	c.concurrentConnectionLimiter <- struct{}{}
	defer func() { <-c.concurrentConnectionLimiter }()

	start := time.Now()
	res, err := c.client.Do(req)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	defer res.Body.Close()

	guessDuration := time.Since(start)

	var guess Guess
	err = json.NewDecoder(res.Body).Decode(&guess)
	if err != nil {
		s.Error = err.Error()
		return s
	}

	if !ValidGuess(guess.Guess, answer) {
		s.Error = fmt.Sprintf("guess is invalid: %s", guess.Guess)
		return s
	}

	result := GetResult(guess.Guess, answer)

	guessResult := GuessResult{
		ID:     id,
		Result: result,
		Guess:  guess.Guess,
	}

	s.GuessResults = append(s.GuessResults, guessResult)
	s.GuessDurationsNS = append(s.GuessDurationsNS, guessDuration.Nanoseconds())
	s.shouts = append(s.shouts, guess.Shout)

	if guess.Guess == answer {
		s.Correct = true
	}

	return s
}

// this struct includes the player's id to give them certainty about who they were
type PlayerMatchResults struct {
	PlayerID string        `json:"player_id,omitempty"`
	Results  MatchSnapshot `json:"results,omitempty"`
}

func (p *Player) BroadcastMatch(m MatchSnapshot) error {

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

func GetDefinition(uri string) (PlayerDefinition, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", uri), nil)
	if err != nil {
		return PlayerDefinition{}, err
	}

	// We should use a very lenient http.Client here since users could be doing anything.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 500 * time.Second,
	}

	res, err := client.Do(req)
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
