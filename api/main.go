package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	port = "8080"
)

func init() {
	flag.StringVar(&port, "port", port, "port to listen for games on")
}

func main() {

	r := gin.Default()
	api := r.Group("/api")
	api.POST("/match", handleStartMatch)

	r.Run(fmt.Sprintf(":%s", port))

}
