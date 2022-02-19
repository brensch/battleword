package battleword

import (
	"log"
	"sync"
	"time"
)

type Match struct {
	Players []*Player `json:"players,omitempty"`
	Games   []*Game   `json:"games,omitempty"`

	Summary *MatchSummary `json:"summary,omitempty"`

	numRounds  int
	numLetters int

	allWords    []string
	commonWords []string
}

type MatchSummary struct {
	Fastest      Fastest      `json:"fastest,omitempty"`
	MostAccurate MostAccurate `json:"most_accurate,omitempty"`
	Loudest      Loudest      `json:"loudest,omitempty"`
	MostCorrect  MostCorrect  `json:"most_correct,omitempty"`

	GamesFastest      MostGames `json:"games_fastest,omitempty"`
	GamesLoudest      MostGames `json:"games_loudest,omitempty"`
	GamesMostAccurate MostGames `json:"games_most_accurate,omitempty"`
}

// InitMatch generates all the games for the match and populates player information and other match level metadata
func InitMatch(allWords, commonWords []string, playerURIs []string, numLetters, numRounds, numGames int) (*Match, error) {

	games := make([]*Game, numGames)
	for i := 0; i < numGames; i++ {
		games[i] = InitGame(commonWords, numLetters, numRounds)
	}

	var wgGenerate, wgListen sync.WaitGroup
	playerCHAN := make(chan *Player)
	var players []*Player

	wgListen.Add(1)
	go func() {
		defer wgListen.Done()
		for player := range playerCHAN {
			players = append(players, player)

		}
	}()

	for _, playerURI := range playerURIs {
		wgGenerate.Add(1)
		go func(playerURI string) {
			defer wgGenerate.Done()
			player, err := InitPlayer(playerURI)
			if err != nil {
				log.Printf("player %s failed to respond: %+v", playerURI, err)
				return
			}

			playerCHAN <- player

		}(playerURI)
	}

	wgGenerate.Wait()
	close(playerCHAN)
	wgListen.Wait()

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
	player *Player
}

// PlayGame takes one of the games for a match and sends it to all players.
// as players finish their games, they are sent back on a channel to be summarised.
// overall results are calculated as each players individual results arrive to
// avoid having an extra loop through all player results once they're all finished.
func (m *Match) PlayGame(g *Game) {

	var wgGames, wgResults sync.WaitGroup
	playerResultsCHAN := make(chan playerResult)

	summary := &GameSummary{
		Start: time.Now(),
	}

	fastestTime := 100 * time.Hour
	bestAccuracy := m.numRounds + 1
	// this tracks if it's a tie since it's easy for multiple players to guess in the same number of rounds
	playersWithSameBestAccuracy := 0

	// listen for the summary
	wgResults.Add(1)
	go func() {
		defer wgResults.Done()
		for result := range playerResultsCHAN {

			// append game to player here since it's serialised across goroutines
			result.player.mu.Lock()
			result.player.Games = append(result.player.Games, result.state)
			result.player.mu.Unlock()

			// then perform all calculations across all the players for this particular game
			if result.state.TotalTime < fastestTime {
				summary.Fastest = Fastest{
					PlayerID: result.player.ID,
					Time:     result.state.TotalTime,
				}
				fastestTime = result.state.TotalTime
			}

			numGuesses := len(result.state.Guesses)
			if !result.state.Correct {
				numGuesses++
			}
			if numGuesses == bestAccuracy {
				playersWithSameBestAccuracy++
			}

			if numGuesses < bestAccuracy {
				summary.MostAccurate = MostAccurate{
					PlayerID:           result.player.ID,
					AverageGuessLength: float64(numGuesses),
				}
				bestAccuracy = numGuesses
				playersWithSameBestAccuracy = 0
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
				player: player,
			}
		}(player)
	}

	// wait for games and summary
	wgGames.Wait()
	close(playerResultsCHAN)
	wgResults.Wait()
	summary.End = time.Now()

	if playersWithSameBestAccuracy > 0 {
		summary.MostAccurate = MostAccurate{}
	}

	g.Summary = summary
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

// Summarise aggregates statistics across the whole match
func (m *Match) Summarise() {

	fastest := Fastest{Time: 100 * time.Hour}
	mostAccurate := MostAccurate{AverageGuessLength: 11}
	var mostCorrect MostCorrect
	var loudest Loudest

	for _, player := range m.Players {
		player.Summarise()
		if player.Summary.Disqualified {
			continue
		}

		if player.Summary.AverageGuesses < mostAccurate.AverageGuessLength {
			mostAccurate.AverageGuessLength = player.Summary.AverageGuesses
			mostAccurate.PlayerID = player.ID
		}

		if player.Summary.TotalTime < fastest.Time {
			fastest.Time = player.Summary.TotalTime
			fastest.PlayerID = player.ID
		}

		if player.Summary.GamesWon > mostCorrect.CorrectGames {
			mostCorrect.CorrectGames = player.Summary.GamesWon
			mostCorrect.PlayerID = player.ID
		}

		if player.Summary.TotalVolume > loudest.Volume {
			loudest.Volume = player.Summary.TotalVolume
			loudest.PlayerID = player.ID
		}
	}

	gamesFastestTally := make(map[string]int)
	gamesLoudestTally := make(map[string]int)
	gamesMostAccurateTally := make(map[string]int)

	for _, game := range m.Games {
		gamesFastestTally[game.Summary.Fastest.PlayerID]++
		gamesLoudestTally[game.Summary.Loudest.PlayerID]++
		gamesMostAccurateTally[game.Summary.MostAccurate.PlayerID]++
	}

	var gamesFastest, gamesLoudest, gamesMostAccurate MostGames

	for playerID, gameCount := range gamesFastestTally {
		if gameCount > gamesFastest.Count {
			gamesFastest.PlayerID = playerID
			gamesFastest.Count = gameCount
		}
	}

	for playerID, gameCount := range gamesLoudestTally {
		if gameCount > gamesLoudest.Count && playerID != "" {
			gamesLoudest.PlayerID = playerID
			gamesLoudest.Count = gameCount
		}
	}

	for playerID, gameCount := range gamesMostAccurateTally {
		if gameCount > gamesMostAccurate.Count {
			gamesMostAccurate.PlayerID = playerID
			gamesMostAccurate.Count = gameCount
		}
	}

	m.Summary = &MatchSummary{
		Loudest:      loudest,
		MostAccurate: mostAccurate,
		MostCorrect:  mostCorrect,
		Fastest:      fastest,

		GamesFastest:      gamesFastest,
		GamesLoudest:      gamesLoudest,
		GamesMostAccurate: gamesMostAccurate,
	}
}
