package battleword

import "testing"

func TestInitGame(t *testing.T) {

	commonWords := []string{
		"beast",
		"feast",
	}

	g := InitGame(commonWords, 5, 6)

	if g.Answer == "" {
		t.Log("got empty answer")
		t.Fail()
	}

	if g.ID == "" {
		t.Log("got empty id")
		t.Fail()
	}

}
