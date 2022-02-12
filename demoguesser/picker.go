package main

import (
	"math/rand"
	"time"

	"github.com/brensch/battleword"
)

func GuessWord() string {
	rand.Seed(time.Now().UnixNano())
	return battleword.CommonWords[rand.Intn(len(battleword.CommonWords))]
}
