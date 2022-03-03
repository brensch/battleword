package main

import (
	"net/http"

	"github.com/brensch/battleword"
	"github.com/gin-gonic/gin"
)

type StartMatchRequest struct {
	Letters int      `json:"letters,omitempty"`
	Games   int      `json:"games,omitempty"`
	Players []string `json:"players,omitempty"`
}

type StartMatchResponse struct {
	UUID    string                         `json:"uuid,omitempty"`
	Players []*battleword.PlayerDefinition `json:"players,omitempty"`
}

func (s *apiStore) handleStartMatch(c *gin.Context) {

	var req StartMatchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	match, err := battleword.InitMatch(s.log, battleword.AllWords, battleword.CommonWords, req.Players, req.Letters, 6, req.Games)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// background match calls here. they will handle updating firestore internally.
	go func() {
		match.Start()
		match.Summarise()
		match.Broadcast()
	}()

	var playerDefinitions []*battleword.PlayerDefinition
	for _, player := range match.Players {
		playerDefinitions = append(playerDefinitions, player.Definition)
	}

	c.JSON(200, StartMatchResponse{
		UUID:    match.UUID,
		Players: playerDefinitions,
	})
}
