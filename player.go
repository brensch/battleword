package battleword

import (
	"bytes"
	"context"
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
	ID          string            `json:"player_id,omitempty"`
	Definition  PlayerDefinition  `json:"definition,omitempty"`
	GamesPlayed []PlayerGameState `json:"games_played,omitempty"`

	Finish time.Time `json:"finish,omitempty"`

	connection PlayerConnection

	mu  *sync.Mutex
	log logrus.FieldLogger
}

type PlayerDefinition struct {
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	ConcurrentConnLimit int    `json:"concurrent_connection_limit,omitempty"`
	Colour              string `json:"colour,omitempty"`
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
	ID string `json:"guess_id,omitempty"`

	Start  time.Time `json:"start,omitempty"`
	Finish time.Time `json:"finish,omitempty"`

	Guess  string `json:"guess,omitempty"`
	Result []int  `json:"result,omitempty"`
}

type PlayerGameState struct {
	GameID string `json:"game_id,omitempty"`

	Start  time.Time `json:"start,omitempty"`
	Finish time.Time `json:"finish,omitempty"`

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

func (p *Player) PlayMatch(ctx context.Context, games []Game) {
	log := p.log.WithField("player_id", p.ID)
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

	// This is used so that the games themselves are staggered.
	// Was noticing that when you start all the games at once you don't get any finished until the very end.
	gameStaggerCHAN := make(chan struct{}, p.Definition.ConcurrentConnLimit+2)

	for _, game := range games {
		wgGenerate.Add(1)
		go func(game Game) {
			defer wgGenerate.Done()

			// only start playing the game when there's no more than the concurrent connection limit going on at once
			gameStaggerCHAN <- struct{}{}
			defer func() { <-gameStaggerCHAN }()

			gameStateCHAN <- PlayGame(ctx, log, p.connection, game)
		}(game)
	}

	wgGenerate.Wait()
	close(gameStateCHAN)
	wgListen.Wait()

	p.mu.Lock()
	p.Finish = time.Now()
	p.mu.Unlock()

	log.Info("player finished match")
}

func PlayGame(ctx context.Context, log logrus.FieldLogger, c PlayerConnection, g Game) PlayerGameState {

	state := PlayerGameState{
		GameID: g.ID,
		Start:  time.Now(),
	}
	log = log.WithField("game_id", g.ID)
	log.Info("player started game")

	for {
		if ctx.Err() != nil {
			state.Error = "match was cancelled"
			state.Finish = time.Now()
			return state
		}

		state = GetNextState(ctx, log, c, state, g.Answer)

		// https://i.redd.it/cw0cedsc93h81.jpg
		if state.Correct || state.Error != "" || len(state.GuessResults) == g.numRounds {
			state.Finish = time.Now()
			log.Info("player finished game")
			return state
		}
	}

}

// this struct includes the player's id to give them certainty about who they were
type PlayerMatchResults struct {
	PlayerID string        `json:"player_id,omitempty"`
	Results  MatchSnapshot `json:"results,omitempty"`
}

func (p *Player) BroadcastMatch(m MatchSnapshot) error {

	p.log.Info("Broadcasting match results")

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
