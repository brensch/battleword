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

// get players
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

func (m *Match) PlayGame(g *Game) {

	var wgGames, wgResults sync.WaitGroup
	playerResultsCHAN := make(chan playerResult)

	results := &GameResult{
		Start: time.Now(),
	}

	fastestTime := 100 * time.Hour
	bestAccuracy := m.numRounds + 1

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

	wgGames.Wait()
	close(playerResultsCHAN)

	wgResults.Wait()
	results.End = time.Now()

	g.Result = results
}

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
