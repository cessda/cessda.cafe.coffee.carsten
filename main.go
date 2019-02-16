package main

import (
  "time"
  "net/http"
  "github.com/sirupsen/logrus"
  "github.com/gin-gonic/gin"
  "github.com/atarantini/ginrequestid"
)

func myRequestLogger(log *logrus.Logger) gin.HandlerFunc {
  return func(c *gin.Context) {

    path := c.Request.URL.Path

    start := time.Now()

    log.WithFields(logrus.Fields{
      "method": "GET",
      "requestId": c.MustGet("RequestId"),
      "endpoint": path,
    }).Debug("request received")


    c.Next()

    elapsed := time.Since(start)

    log.WithFields(logrus.Fields{
      "method": "GET",
      "requestId": c.MustGet("RequestId"),
      "status": http.StatusOK,
      "responsetime": elapsed,
      "endpoint": path,
    }).Debug("request answered")

  }
}

func setupRouter() *gin.Engine {
  r := gin.New()

  r.Use(ginrequestid.RequestId())

  log := logrus.New()
  log.SetLevel(logrus.DebugLevel)

  r.Use(myRequestLogger(log), gin.Recovery())

  r.GET("/", func(c *gin.Context) {
    c.String(http.StatusOK, "Hello World!")
  })

  r.GET("/healthcheck", func(c *gin.Context) {

    c.JSON(http.StatusOK, gin.H{"message": "Ok"})

  })

  r.GET("/order-history", func(c *gin.Context) {

    c.JSON(http.StatusOK, getAllOrders())

  })

  r.GET("/order-history/:ID", func(c *gin.Context) {
    ID := c.Param("ID")
    order, err := getOrderbyID(ID)

    if err != nil {
      c.JSON(http.StatusNotFound, err)
    } else {
      c.JSON(http.StatusOK, order)
    }
  })

  return r
}

func main() {

  r := setupRouter()
  r.Run(":1337")

}
