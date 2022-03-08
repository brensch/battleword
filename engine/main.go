package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brensch/battleword"
	"github.com/sirupsen/logrus"
)

var (
	NumRounds  = 6
	NumLetters = 5
	NumGames   = 3

	PlayerURIsJoined = "http://localhost:8080"
)

func init() {
	flag.IntVar(&NumRounds, "num_rounds", NumRounds, "the number of rounds you want in each game")
	flag.IntVar(&NumLetters, "num_letters", NumLetters, "the number of letters the word should be")
	flag.IntVar(&NumGames, "num_games", NumGames, "how many games to play in the match")

	flag.StringVar(&PlayerURIsJoined, "apis", PlayerURIsJoined, "the location of all players' apis in the same order as their names, separated by commas")

}

func main() {

	flag.Parse()

	log := logrus.New()

	log.Info("started game")

	playerURIs := strings.Split(PlayerURIsJoined, ",")

	if playerURIs[0] == "" {
		log.Println("you need to define player endpoints")
		return
	}

	filename := fmt.Sprintf("results-%s.json", time.Now().Format(time.RFC3339))
	f, err := os.Create(filename)
	if err != nil {
		log.Println("couldn't create file", filename)
		return
	}
	defer f.Close()

	match, err := battleword.InitMatch(log, battleword.AllWords, battleword.CommonWords, playerURIs, NumLetters, NumRounds, NumGames)
	if err != nil {
		log.Println("got error initing game", err)
		return
	}

	match.Start()
	match.Broadcast()

	log.Println("game finished")
	// gameJSON, _ := json.Marshal(match)

	err = json.NewEncoder(f).Encode(match)
	if err != nil {
		log.Println("couldn't write to file", err)
		return
	}

	log.Println("final result saved to file", filename)

}
