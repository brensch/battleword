package battleword

import (
	"testing"
)

func TestInitGameState(t *testing.T) {
	players := []*Player{
		{
			Name: "brend",
		},
	}
	g, err := InitGameState(players, 5, 6)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(g.Players) != len(players) {
		t.Log("mismatched player length")
		t.FailNow()
	}
}
