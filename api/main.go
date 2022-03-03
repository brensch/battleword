package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	port = "8080"
)

func init() {
	flag.StringVar(&port, "port", port, "port to listen for games on")
}

type apiStore struct {
	log logrus.FieldLogger
}

func main() {

	log := logrus.New()
	// log.SetLevel(logrus.DebugLevel)
	s := &apiStore{
		log: log,
	}

	r := gin.New()
	r.Use(MiddlewareLogger(log))
	api := r.Group("/api")
	api.POST("/match", s.handleStartMatch)

	r.Run(fmt.Sprintf(":%s", port))

}

func MiddlewareLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := "api call"
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
