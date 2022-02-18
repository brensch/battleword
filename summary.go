package battleword

import "time"

type Fastest struct {
	PlayerID string        `json:"player_id,omitempty"`
	Time     time.Duration `json:"time,omitempty"`
}

type MostAccurate struct {
	PlayerID           string  `json:"player_id,omitempty"`
	AverageGuessLength float64 `json:"average_guess_length,omitempty"`
}

type Loudest struct {
	PlayerID string  `json:"player_id,omitempty"`
	Volume   float64 `json:"volume,omitempty"`
}

type MostCorrect struct {
	PlayerID     string `json:"player_id,omitempty"`
	CorrectGames int    `json:"correct_games,omitempty"`
}

type MostGames struct {
	PlayerID string `json:"player_id,omitempty"`
	Count    int    `json:"count,omitempty"`
}
