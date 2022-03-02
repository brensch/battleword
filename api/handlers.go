package main

import (
	"log"
	"net/http"

	"github.com/brensch/battleword"
	"github.com/gin-gonic/gin"
)

type StartMatchRequest struct {
	Letters int      `json:"letters,omitempty"`
	Games   int      `json:"games,omitempty"`
	Players []string `json:"players,omitempty"`
}

func handleStartMatch(c *gin.Context) {

	var req StartMatchRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	match, err := battleword.InitMatch(battleword.AllWords, battleword.CommonWords, req.Players, req.Letters, 6, req.Games)
	if err != nil {
		log.Println("got error initing game", err)
		return
	}

	// TODO: obviously needs to be backgrounded, written to firestore etc.
	match.Start()
	match.Summarise()
	match.Broadcast()

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
