package battleword

import (
	"log"
	"sync"
	"time"
)

type Match struct {
	Players []*Player `json:"players,omitempty"`
	Games   []*Game   `json:"games,omitempty"`

	numRounds  int
	numLetters int

	allWords    []string
	commonWords []string
}

// InitMatch generates all the games for the match and populates player information and other match level metadata
func InitMatch(allWords, commonWords []string, players []*Player, numLetters, numRounds, numGames int) (*Match, error) {

	games := make([]*Game, numGames)
	for i := 0; i < numGames; i++ {
		games[i] = InitGame(commonWords, numLetters, numRounds)
	}

	return &Match{
		Players: players,
		Games:   games,

		numLetters: numLetters,
		numRounds:  numRounds,

		allWords:    allWords,
		commonWords: commonWords,
	}, nil

}

// Start kicks off all the games as goroutines and waits for them to complete
func (m *Match) Start() {

	var wg sync.WaitGroup

	for _, game := range m.Games {
		wg.Add(1)
		go func(game *Game) {
			defer wg.Done()
			m.PlayGame(game)
		}(game)
	}

	wg.Wait()
}

type playerResult struct {
	state  *PlayerGameState
	player PlayerDefinition
}

// PlayGame takes one of the games for a match and sends it to all players.
// as players finish their games, they are sent back on a channel to be summarised.
// overall results are calculated as each players individual results arrive to
// avoid having an extra loop through all player results once they're all finished.
func (m *Match) PlayGame(g *Game) {

	var wgGames, wgResults sync.WaitGroup
	playerResultsCHAN := make(chan playerResult)

	results := &GameResult{
		Start: time.Now(),
	}

	fastestTime := 100 * time.Hour
	bestAccuracy := m.numRounds + 1

	// listen for the results
	wgResults.Add(1)
	go func() {
		defer wgResults.Done()
		for result := range playerResultsCHAN {

			if result.state.TotalTime < fastestTime {
				results.Fastest = FastestPlayer{
					Player: result.player,
					Time:   result.state.TotalTime,
				}
				fastestTime = result.state.TotalTime
			}

			if len(result.state.Guesses) < bestAccuracy {
				results.MostAccurate = MostAccuratePlayer{
					Player:             result.player,
					AverageGuessLength: float64(len(result.state.Guesses)),
				}
				bestAccuracy = len(result.state.Guesses)
			}

		}
	}()

	// play the games
	for _, player := range m.Players {
		wgGames.Add(1)
		go func(player *Player) {
			defer wgGames.Done()
			state := player.PlayGame(g)
			playerResultsCHAN <- playerResult{
				state:  state,
				player: player.Definition,
			}
		}(player)
	}

	// wait for games and results
	wgGames.Wait()
	close(playerResultsCHAN)
	wgResults.Wait()
	results.End = time.Now()

	g.Result = results
}

// Broadcast takes the results of the match and sends it to all players
func (m *Match) Broadcast() {

	var wg sync.WaitGroup

	for _, player := range m.Players {
		wg.Add(1)
		go func(player *Player) {
			defer wg.Done()
			err := player.BroadcastMatch(m)
			if err != nil {
				log.Println(err)
			}
		}(player)
	}

	wg.Wait()
}
