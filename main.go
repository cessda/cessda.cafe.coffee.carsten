package main

import (
  "time"
  "net/http"
  "github.com/sirupsen/logrus"
  "github.com/toorop/gin-logrus"
  "github.com/gin-gonic/gin"
  "github.com/atarantini/ginrequestid"
)

func setupRouter() *gin.Engine {
  r := gin.New()

  r.Use(ginrequestid.RequestId())

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
      "requestId": c.MustGet("RequestId"),
      "endpoint": "healthcheck",
    }).Debug("request received")

    c.JSON(http.StatusOK, gin.H{"message": "Ok"})

    elapsed := time.Since(start)

    log.WithFields(logrus.Fields{
      "method": "GET",
      "requestId": c.MustGet("RequestId"),
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
