package battleword

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMatchInit(t *testing.T) {

	playerURIs := []string{
		"http://localhost:8080",
	}

	// players := []*Player{
	// 	{
	// 		Definition: &PlayerDefinition{"", "brend", ""},
	// 		connection: &PlayerConnection{
	// 			uri:    "http://localhost:8080",
	// 			client: http.DefaultClient,
	// 		},
	// 	}}

	match, err := InitMatch(logrus.New(), AllWords, CommonWords, playerURIs, 5, 6, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, game := range match.Games {
		fmt.Println(game.Answer)
	}
}

func TestMatchStart(t *testing.T) {

	playerURIs := []string{
		"http://localhost:8080",
	}
	// players := []*Player{
	// 	{
	// 		Definition: &PlayerDefinition{"", "brend", ""},
	// 		connection: &PlayerConnection{
	// 			uri:    "http://localhost:8080",
	// 			client: http.DefaultClient,
	// 		},
	// 	},
	// }

	match, err := InitMatch(logrus.New(), AllWords, CommonWords, playerURIs, 5, 6, 10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	match.Start()
}
