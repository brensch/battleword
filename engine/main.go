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

	PlayerURIsJoined = "http://localhost:8080"
)

func init() {
	flag.IntVar(&NumRounds, "num_rounds", NumRounds, "the number of rounds you want in each game")
	flag.IntVar(&NumLetters, "num_letters", NumLetters, "the number of letters the word should be")

	flag.StringVar(&PlayerURIsJoined, "apis", PlayerURIsJoined, "the location of all players' apis in the same order as their names, separated by commas")

}

func main() {

	flag.Parse()

	log.Printf("started game")

	playerURIs := strings.Split(PlayerURIsJoined, ",")

	if playerURIs[0] == "" {
		log.Println("you need to define player endpoints")
		return
	}

	// for i, name := range playerNames {
	// 	battleword.InitPlayer(playerURIs[i])
	// 	players = append(players, battleword.InitPlayer(playerURIs[i]))

	// }

	match, err := battleword.InitMatch(battleword.AllWords, battleword.CommonWords, playerURIs, 5, 6, 100)
	if err != nil {
		log.Println("got error initing game", err)
		return
	}

	match.Start()
	match.Summarise()
	match.Broadcast()

	log.Println("game finished")
	gameJSON, _ := json.Marshal(match)
	log.Println("final result:", string(gameJSON))

}
