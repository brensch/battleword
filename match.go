package battleword

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Match struct {
	uuid string

	players []*Player
	games   []Game

	start time.Time
	end   time.Time

	numRounds  int
	numLetters int

	allWords    []string
	commonWords []string

	// this is used to stop writing to allow us to take a snapshot for upload
	mu  *sync.Mutex
	log logrus.FieldLogger
}

// MatchSnapshot is what clients can use to get the current state of the game,
// and what gets sent to the contestants at the end of matches
type MatchSnapshot struct {
	UUID string `json:"match_id,omitempty"`

	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`

	Players []Player `json:"players,omitempty"`
	Games   []Game   `json:"games,omitempty"`

	RoundsPerGame  int `json:"rounds_per_game,omitempty"`
	LettersPerWord int `json:"letters_per_word,omitempty"`
}

// InitMatch generates all the games for the match and populates player information and other match level metadata
func InitMatch(log logrus.FieldLogger, allWords, commonWords []string, playerURIs []string, numLetters, numRounds, numGames int) (*Match, error) {
	id := uuid.NewString()
	log = log.WithField("match_id", id)
	mu := &sync.Mutex{}

	// this could be a large number so preallocate for speed
	games := make([]Game, numGames)
	for i := 0; i < numGames; i++ {
		games[i] = InitGame(commonWords, numLetters, numRounds)
	}

	var wgGenerate, wgListen sync.WaitGroup
	playerCHAN := make(chan *Player)
	errCHAN := make(chan error)
	var players []*Player
	var errors []error

	wgListen.Add(1)
	go func() {
		defer wgListen.Done()
		for player := range playerCHAN {
			log.
				WithFields(logrus.Fields{
					"player_definition": player.Definition,
				}).
				Debug("got player info")

			players = append(players, player)

		}
	}()

	wgListen.Add(1)
	go func() {
		defer wgListen.Done()
		for err := range errCHAN {
			errors = append(errors, err)
		}
	}()

	for _, playerURI := range playerURIs {
		wgGenerate.Add(1)
		go func(playerURI string) {
			defer wgGenerate.Done()
			player, err := InitPlayer(mu, log, playerURI)
			if err != nil {
				log.
					WithFields(logrus.Fields{
						"uri": playerURI,
					}).
					WithError(err).
					Error("player failed to respond")
				errCHAN <- err
				return
			}

			playerCHAN <- player

		}(playerURI)
	}

	wgGenerate.Wait()
	close(playerCHAN)
	close(errCHAN)
	wgListen.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to contact %d players: %+v", len(errors), errors)
	}

	return &Match{
		uuid: id,

		players: players,
		games:   games,

		numLetters: numLetters,
		numRounds:  numRounds,

		allWords:    allWords,
		commonWords: commonWords,

		mu:  mu,
		log: log,
	}, nil

}

// Start kicks off all the games as goroutines and waits for them to complete
func (m *Match) Start(ctx context.Context) {

	m.log.Info("match started")
	m.start = time.Now()

	var wg sync.WaitGroup

	for _, player := range m.players {
		wg.Add(1)
		go func(player *Player) {
			defer wg.Done()
			player.PlayMatch(ctx, m.games)
		}(player)
	}

	wg.Wait()
	m.end = time.Now()
	m.log.Info("match finished")

}

func (m *Match) Snapshot() MatchSnapshot {

	m.mu.Lock()
	defer m.mu.Unlock()

	// JSON isolator
	body, err := json.Marshal(m.players)
	if err != nil {
		m.log.WithError(err).Error("wtf")
		return MatchSnapshot{}
	}

	var decoupledPlayers []Player
	err = json.Unmarshal(body, &decoupledPlayers)
	if err != nil {
		m.log.Error("wtf")
		return MatchSnapshot{}
	}

	return MatchSnapshot{
		UUID:           m.uuid,
		Games:          m.games,
		Players:        decoupledPlayers,
		RoundsPerGame:  m.numRounds,
		LettersPerWord: m.numLetters,

		Start: m.start,
		End:   m.end,
	}

}

// Broadcast takes the results of the match and sends it to all players
func (m *Match) Broadcast() {

	var wg sync.WaitGroup
	snapshot := m.Snapshot()

	for _, player := range m.players {
		wg.Add(1)
		go func(player *Player) {
			defer wg.Done()
			err := player.BroadcastMatch(snapshot)
			if err != nil {
				m.log.WithError(err).Error("failed to broadcast")
			}
		}(player)
	}

	wg.Wait()
}
