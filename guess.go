package battleword

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Guess struct {
	Guess string `json:"guess,omitempty"`

	// For the lols:
	Shout string `json:"shout,omitempty"`
}

func GetNextState(ctx context.Context, log logrus.FieldLogger, c PlayerConnection, s PlayerGameState, answer string) PlayerGameState {
	id := uuid.New().String()
	log = log.WithField("guess_id", id)
	log.WithFields(logrus.Fields{
		"turn": len(s.GuessResults),
	}).
		Debug("queued getting next playergamestate")

	guessesJson, err := json.Marshal(s)
	if err != nil {
		log.WithError(err).Error("failed to encode playergamestate")
		s.Error = err.Error()
		return s
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/guess", c.uri), bytes.NewReader(guessesJson))
	if err != nil {
		log.WithError(err).Error("failed to form request")
		s.Error = err.Error()
		return s
	}

	req.Header.Add(GuessIDHeader, id)

	// Make sure we don't go over the concurrent connection limit for this player.
	// This will halt until there is a free slot in the concurrent request queue.
	c.concurrentConnectionLimiter <- struct{}{}
	defer func() { <-c.concurrentConnectionLimiter }()

	log.Debug("started getting next playergamestate")

	start := time.Now()
	res, err := c.client.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to get player guess")
		s.Error = err.Error()
		return s
	}
	defer res.Body.Close()
	finish := time.Now()
	guessDuration := time.Since(start)

	// want to get full bytes of what player sent to help them out
	guessBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("failed to read from player response")
		s.Error = err.Error()
		return s
	}

	var guess Guess
	err = json.Unmarshal(guessBytes, &guess)
	if err != nil {
		log.WithError(err).Error("failed to decode player response")
		s.Error = err.Error()
		return s
	}

	if !ValidGuess(guess.Guess, answer) {
		log.Warn("received invalid guess")
		s.Error = fmt.Sprintf("guess is invalid: %s", guess.Guess)
		return s
	}

	result := GetResult(guess.Guess, answer)

	guessResult := GuessResult{
		ID:     id,
		Result: result,
		Guess:  guess.Guess,

		Start:  start,
		Finish: finish,
	}

	s.GuessResults = append(s.GuessResults, guessResult)
	s.GuessDurationsNS = append(s.GuessDurationsNS, guessDuration.Nanoseconds())
	s.shouts = append(s.shouts, guess.Shout)

	if guess.Guess == answer {
		s.Correct = true
	}

	log.Debug("successfully got next playergamestate")

	return s
}
