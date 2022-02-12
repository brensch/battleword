package main

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	"github.com/brensch/battleword"
)

var (
	NumRounds  = 6
	NumLetters = 5

	PlayerNamesJoined = "solvo"
	PlayerURIsJoined  = "http://localhost:8080"
)

func init() {
	flag.IntVar(&NumRounds, "num_rounds", NumRounds, "the number of rounds you want in each game")
	flag.IntVar(&NumLetters, "num_letters", NumLetters, "the number of letters the word should be")

	flag.StringVar(&PlayerNamesJoined, "names", PlayerNamesJoined, "the names of all players, separated by commas")
	flag.StringVar(&PlayerURIsJoined, "apis", PlayerURIsJoined, "the location of all players' apis in the same order as their names, separated by commas")

}

func main() {

	flag.Parse()

	log.Printf("started game, players: %s", PlayerNamesJoined)

	playerNames := strings.Split(PlayerNamesJoined, ",")
	playerURIs := strings.Split(PlayerURIsJoined, ",")

	if playerNames[0] == "" {
		log.Println("you need to define players")
		return
	}

	if len(playerNames) != len(playerURIs) {
		log.Println("you need the same number of names as api locations")
		return
	}

	var players []*battleword.Player

	for i, name := range playerNames {
		players = append(players, battleword.InitPlayer(name, playerURIs[i]))

	}

	gamesState, err := battleword.InitGameState(battleword.AllWords, battleword.CommonWords, players, 5, 6)
	if err != nil {
		log.Println("got error initing game", err)
		return
	}

	gamesState.PlayGame()

	log.Println("game finished")
	gameJSON, _ := json.Marshal(gamesState)
	log.Println("final result:", string(gameJSON))

}
