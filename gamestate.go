package battleword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type GameState struct {
	Players []*Player `json:"players,omitempty"`
	Answer  string    `json:"answer,omitempty"`

	numLetters int
	numRounds  int

	allWords    []string
	commonWords []string
}

func InitGameState(allWords, commonWords []string, players []*Player, numLetters, numRounds int) (*GameState, error) {

	if len(players) == 0 {
		return nil, fmt.Errorf("no players")
	}
	state := &GameState{
		Players: players,
		Answer:  GetRandomWord(commonWords),

		numLetters: numLetters,
		numRounds:  numRounds,

		allWords:    allWords,
		commonWords: commonWords,
	}

	return state, nil
}

func GetRandomWord(words []string) string {
	// don't need to do this any more legit than this i don't think
	rand.Seed(time.Now().UnixNano())

	return words[rand.Intn(len(words))]
}

func (g *GameState) PlayGame() {

	var wg sync.WaitGroup

	for _, player := range g.Players {
		wg.Add(1)
		go func(player *Player) {
			defer wg.Done()
			player.PlayGame(g.Answer, g.numRounds)
		}(player)
	}

	wg.Wait()

	g.BroadcastResults()
}

func (g *GameState) BroadcastResults() {

	gameStateJSON, err := json.Marshal(g)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, player := range g.Players {
		wg.Add(1)

		go func(player *Player) {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/results", player.uri), bytes.NewReader(gameStateJSON))
			if err != nil {
				log.Println(err)
				return
			}

			// TODO: make this a proper client
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			res.Body.Close()
		}(player)
	}

	wg.Wait()

}
