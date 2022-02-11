package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brensch/battleword"
)

func main() {
	http.HandleFunc("/guess", DoGuess)
	http.HandleFunc("/results", ReceiveResults)
	http.ListenAndServe(":8080", nil)
}

func DoGuess(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var prevGuesses battleword.Guesses
	err := json.NewDecoder(r.Body).Decode(&prevGuesses)
	if err != nil {
		log.Println(err)
		return
	}

	prevGuessesJSON, _ := json.Marshal(prevGuesses)

	log.Println(string(prevGuessesJSON))

	guess := battleword.Guess{
		Guess: "beast",
		Shout: "it's not beast",
	}

	// time.Sleep(1 * time.Second)

	err = json.NewEncoder(w).Encode(guess)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveResults(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var finalState battleword.GameState
	err := json.NewDecoder(r.Body).Decode(&finalState)
	if err != nil {
		log.Println(err)
		return
	}

	finalStateJSON, _ := json.Marshal(finalState)

	log.Println(string(finalStateJSON))

}
