package battleword

import (
	"fmt"
	"testing"
)

func TestMatchInit(t *testing.T) {

	players := []*Player{
		InitPlayer("brendan", "a cool guy", "http://localhost:8080"),
	}

	match, err := InitMatch(AllWords, CommonWords, players, 5, 6, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, game := range match.Games {
		fmt.Println(game.Answer)
	}
}

func TestMatchStart(t *testing.T) {

	players := []*Player{
		InitPlayer("brendan", "a cool guy", "http://localhost:8080"),
	}

	match, err := InitMatch(AllWords, CommonWords, players, 5, 6, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	match.Start()
}
