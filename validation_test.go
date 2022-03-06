package battleword

import "testing"

func TestValidGuess(t *testing.T) {
	guess := "beast"
	answer := "beast"
	valid := ValidGuess(guess, answer)
	if !valid {
		t.FailNow()
	}
}

type ResultTestCase struct {
	guess  string
	answer string
	result []int
}

var (
	resultTestCases = []ResultTestCase{
		{"beast", "beast", []int{2, 2, 2, 2, 2}},
		{"least", "beast", []int{0, 2, 2, 2, 2}},
		{"beast", "beasy", []int{2, 2, 2, 2, 0}},
		{"trees", "beast", []int{1, 0, 1, 0, 1}},
		{"trees", "ulcer", []int{0, 1, 0, 2, 0}},
		{"ruler", "ulcer", []int{0, 1, 1, 2, 2}},
		{"bluer", "ulcer", []int{0, 2, 1, 2, 2}},
		{"seers", "ulcer", []int{0, 1, 0, 1, 0}},
		{"seers", "ulcee", []int{0, 1, 1, 0, 0}},
		{"seerse", "elceec", []int{0, 1, 1, 0, 0, 1}},
	}
)

func TestGetResults(t *testing.T) {
	for _, testCase := range resultTestCases {
		result := GetResult(testCase.guess, testCase.answer)

		// compare each character in the result
		for i := 0; i < len(result.Result); i++ {
			if result.Result[i] != testCase.result[i] {
				t.Logf(
					"got mismatch. guess: %s, answer %s, result %+v, expected result %+v",
					testCase.guess,
					testCase.answer,
					result,
					testCase.result,
				)
				t.FailNow()
			}
		}
	}
}

// BenchmarkGetResults tests iterations of all test cases.
// may want to add one benchmark for each test in the future.
func BenchmarkGetResults(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, testCase := range resultTestCases {
			GetResult(testCase.guess, testCase.answer)
		}
	}
}
