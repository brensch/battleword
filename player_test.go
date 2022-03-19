package battleword

import (
	"fmt"
	"testing"
)

var (
	validColours = []string{
		"#60283c",
		"#60283C",
		"#4b2860",
	}
	invalidColours = []string{
		"60283c",
		"#60283C1",
		"##4b2860",
	}
)

func TestValidColour(t *testing.T) {

	for _, validColour := range validColours {

		valid, err := ValidColour(validColour)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		if !valid {
			t.Log("should be valid")
			t.FailNow()
		}
	}

	for _, invalidColour := range invalidColours {

		valid, _ := ValidColour(invalidColour)
		if valid {
			t.Log("should not be valid")
			t.FailNow()
		}
	}

}

func TestColourFromString(t *testing.T) {

	colour := ColourFromString("yolo swagginss")
	fmt.Println(colour)

	valid, err := ValidColour(colour)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !valid {
		t.Log("did not get valid colour from colourfromstring")
		t.Fail()
	}
}

func FuzzColourFromString(f *testing.F) {
	f.Add("yeet")
	f.Fuzz(func(t *testing.T, s string) {
		colour := ColourFromString(s)
		fmt.Println(colour)

		valid, err := ValidColour(colour)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		if !valid {
			t.Log("did not get valid colour from colourfromstring")
			t.Fail()
		}
	})
}
