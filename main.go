package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

  r.GET("/healthcheck", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Ok"})
	})

	return r
}

func main() {

	r := setupRouter()
	r.Run(":1337")

}
