package main

import (
  "time"
  "net/http"
  "github.com/sirupsen/logrus"
  "github.com/gin-gonic/gin"
  "github.com/atarantini/ginrequestid"
)

type coffeeOrder struct {
  Type      string `json:"type"`
}


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

  r.GET("/status", func(c *gin.Context) {

    if systemBrewing() {
      c.JSON(http.StatusConflict, gin.H{"message": "System Brewing -- Please wait!"})
    } else if orderWaiting() {
      c.JSON(http.StatusUnauthorized, gin.H{"message": "Coffee Waiting -- Come and get it!"})
    } else {
      c.JSON(http.StatusOK, gin.H{"message": "System Ready"})
    }
  })

  r.GET("/order-history", func(c *gin.Context) {

    c.JSON(http.StatusOK, getAllOrders())

  })

  r.GET("/order-history/:ID", func(c *gin.Context) {
    ID := c.Param("ID")
    order, success := getOrderbyID(ID)

    if ! success {
      c.JSON(http.StatusNotFound, gin.H{"message": "Order not found!"})
    } else {
      c.JSON(http.StatusOK, order)
    }
  })

  r.POST("/order-coffee", func(c *gin.Context) {
    var incomingOrder coffeeOrder
    c.BindJSON(&incomingOrder)

    result, success := newOrder(incomingOrder.Type)

    if !success {
      log.WithFields(logrus.Fields{
        "requestId": c.MustGet("RequestId"),
      }).Debug("Error placing order")
      c.JSON(http.StatusBadRequest, gin.H{"message": "Error!"})
    } else {
      log.WithFields(logrus.Fields{
        "requestId": c.MustGet("RequestId"),
        "type": result.Type,
        "readyAt": result.ReadyAt,
      }).Debug("Order placed")
      c.JSON(http.StatusOK, result)
    }

  })

  r.GET("/retrieve-coffee/:ID", func(c *gin.Context) {
    ID := c.Param("ID")
    order, success := retrieveOrder(ID)
      log.WithFields(logrus.Fields{
        "requestId": c.MustGet("RequestId"),
        "success": success,
      }).Debug("Request to retrieve Coffee")
    if ! success {
      log.WithFields(logrus.Fields{
        "requestId": c.MustGet("RequestId"),
        "success": success,
        "order-id": ID,
      }).Debug("Failed to retrieve Coffee")
      c.JSON(http.StatusBadRequest, gin.H{"message": "Nope!"})
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
