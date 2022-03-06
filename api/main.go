package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

const (
	// cloud run specific env vars
	// the vendor lock in is real
	EnvVarService       = "K_SERVICE"
	EnvVarRevision      = "K_REVISION"
	EnvVarConfiguration = "K_CONFIGURATION"
	EnvVarPort          = "PORT"

	defaultPort = "8080"

	FirestoreMatchCollection = "matches"
)

type apiStore struct {
	log      logrus.FieldLogger
	fsClient *firestore.Client
}

func main() {

	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "battleword"}
	opt := option.WithCredentialsFile("key.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
	}

	fsClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer fsClient.Close()

	// using env var since cloud run uses it
	port := os.Getenv(EnvVarPort)
	if port == "" {
		port = defaultPort
	}

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339Nano,
	})
	s := &apiStore{
		log:      log,
		fsClient: fsClient,
	}

	service := os.Getenv(EnvVarService)
	revision := os.Getenv(EnvVarRevision)
	configuration := os.Getenv(EnvVarConfiguration)

	// service being set indicates we're running on cloud run
	if service != "" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(logrus.InfoLevel)
	}

	r := gin.New()
	r.Use(cors.Default())

	r.Use(MiddlewareLogger(log))
	api := r.Group("/api")
	api.POST("/match", s.handleStartMatch)

	log.
		WithFields(logrus.Fields{
			"service":       service,
			"revision":      revision,
			"configuration": configuration,
		}).
		Info("app starting")

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
			"hostname":    hostname,
			"status_code": statusCode,
			"latency":     latency, // time to process
			"client_ip":   clientIP,
			"method":      c.Request.Method,
			"path":        path,
			"referer":     referer,
			"data_length": dataLength,
			"user_agent":  clientUserAgent,
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
