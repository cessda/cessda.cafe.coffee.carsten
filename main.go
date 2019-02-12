package main

import (
  "time"
	"net/http"
  "github.com/sirupsen/logrus"
  "github.com/toorop/gin-logrus"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

  log := logrus.New()
  log.SetLevel(logrus.DebugLevel)
  //log.SetFormatter(&log.JSONFormatter{})
  r.Use(ginlogrus.Logger(log), gin.Recovery())

  r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

  r.GET("/healthcheck", func(c *gin.Context) {

    start := time.Now()

    log.WithFields(logrus.Fields{
      "method": "GET",
      "endpoint": "healthcheck",
    }).Debug("request received")

    c.JSON(http.StatusOK, gin.H{"message": "Ok"})

    elapsed := time.Since(start)

    log.WithFields(logrus.Fields{
      "method": "GET",
      "status": http.StatusOK,
      "responsetime": elapsed,
      "endpoint": "healthcheck",
    }).Debug("request answered")

	})

	return r
}

func main() {

	r := setupRouter()
	r.Run(":1337")

}
