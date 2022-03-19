package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/brensch/battleword"
)

var (
	port = "8080"
)

func init() {

	flag.StringVar(&port, "port", port, "port to listen for games on")

}

func main() {

	flag.Parse()

	log.Println("I am howyagoin. I'm a bit how ya goin")
	log.Println("waiting to receive a wordle")
	log.Printf("listening on port %s", port)

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/guess", DoGuess)
	http.HandleFunc("/results", ReceiveResults)
	http.HandleFunc("/ping", DoPing)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}

func GuessWord() string {

	return battleword.CommonWords[rand.Intn(len(battleword.CommonWords))]
}

func RandomShout() string {

	shouts := []string{
		"wordle is fun, but for how long?",
		"you will one day be dust, but i will always be solvo",
		"what's the point of anything?",
		"there has to be a better strat than this",
		"i wonder if a human could respond to the api and compete against machines",
	}

	return shouts[rand.Intn(len(shouts))]
}

func DoGuess(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received guess request ID %s. \n", r.Header.Get(battleword.GuessIDHeader))
	if r.Method != http.MethodPost {
		return
	}

	var prevGuesses battleword.PlayerGameState
	err := json.NewDecoder(r.Body).Decode(&prevGuesses)
	if err != nil {
		log.Println(err)
		return
	}

	words, err := GetPossibleWords(prevGuesses.GuessResults)
	if err != nil {
		log.Println(err)
		return
	}

	// pick a randooo
	word := words[rand.Intn(len(words))]

	guess := battleword.Guess{
		Guess: word,
		Shout: RandomShout(),
	}

	// time.Sleep(100 * time.Millisecond)

	err = json.NewEncoder(w).Encode(guess)
	if err != nil {
		log.Println(err)
		return
	}
	prevGuessesJSON, _ := json.Marshal(prevGuesses)
	log.Printf("Making informed guess for game %s, turn %d: %s\n", r.Header.Get(battleword.GuessIDHeader), len(prevGuesses.GuessResults), word)
	log.Printf("Request ID %s. Body: %s\n", r.Header.Get(battleword.GuessIDHeader), prevGuessesJSON)
}

func GetPossibleWords(prevGuessResults []battleword.GuessResult) ([]string, error) {

	possibleWords := battleword.CommonWords
	for _, prevGuessResult := range prevGuessResults {
		var newPossibleWords []string
		for _, newGuess := range possibleWords {
			if WordPossible(newGuess, prevGuessResult) {
				newPossibleWords = append(newPossibleWords, newGuess)
			}
		}
		possibleWords = newPossibleWords
	}
	return possibleWords, nil
}

// Don't steal this boi
func WordPossible(newGuess string, prevGuessResult battleword.GuessResult) bool {

	newResult := battleword.GetResult(prevGuessResult.Guess, newGuess)
	for i := 0; i < len(newGuess); i++ {
		if newResult[i] != prevGuessResult.Result[i] {
			return false
		}
	}

	return true
}

func DoPing(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		return
	}

	log.Println("received ping")

	definition := &battleword.PlayerDefinition{
		Name:        "howyagoin",
		Description: "the average",
	}

	err := json.NewEncoder(w).Encode(definition)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveResults(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var finalState battleword.PlayerMatchResults
	err := json.NewDecoder(r.Body).Decode(&finalState)
	if err != nil {
		log.Println(err)
		return
	}

	var us battleword.Player
	found := false
	for _, player := range finalState.Results.Players {
		if player.ID == finalState.PlayerID {
			us = player
			found = true
		}
	}

	if !found {
		log.Println("We weren't in the results. strange")
		return
	}

	gamesWon := 0
	for _, game := range us.GamesPlayed {
		if game.Correct {
			gamesWon++
		}
	}

	finalStateJSON, _ := json.Marshal(finalState)

	log.Printf("Match %s concluded, we got %d words right. Body: %s", finalState.Results.UUID, gamesWon, finalStateJSON)

}
