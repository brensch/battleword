package battleword

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestValidDefinition(t *testing.T) {
	defString := []byte(`
	{
		"name": "Sendooooo",
		"description": "yoloest",
		"concurrent_connection_limit": 3,
		"colour": "#fcba03"
	}
	`)

	var def PlayerDefinition
	err := json.Unmarshal(defString, &def)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	fmt.Println(def.Colour)

	valid, err := ValidDefinition(def)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !valid {
		t.Log("not valid")
		t.FailNow()
	}

}
