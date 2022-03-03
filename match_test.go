package battleword

import (
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	gameCount = 10
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

	match, err := InitMatch(logrus.New(), AllWords, CommonWords, playerURIs, 5, 6, gameCount)
	if err != nil {
		t.Log(err)
		t.Skip()
		return
	}

	if len(match.Games) != gameCount {
		t.Log("wrong game count", len(match.Games), gameCount)
		t.FailNow()
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

	match, err := InitMatch(logrus.New(), AllWords, CommonWords, playerURIs, 5, 6, gameCount)
	if err != nil {
		t.Log(err)
		t.Skip()
		return
	}

	match.Start()
}
