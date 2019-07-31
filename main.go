package main

import (
	"github.com/atarantini/ginrequestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type coffeeJob struct {
	Product string `json:"product"`
}

func myRequestLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path

		start := time.Now()

		log.WithFields(logrus.Fields{
			"method":    "GET",
			"requestId": c.MustGet("RequestId"),
			"endpoint":  path,
		}).Debug("request received")

		c.Next()

		elapsed := time.Since(start)

		log.WithFields(logrus.Fields{
			"method":       "GET",
			"requestId":    c.MustGet("RequestId"),
			"status":       http.StatusOK,
			"responsetime": elapsed,
			"endpoint":     path,
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

		systemStatusCode, systemStatusMessage := systemStatus()

		c.JSON(systemStatusCode, gin.H{"message": systemStatusMessage})

	})

	r.GET("/job-history", func(c *gin.Context) {

		c.JSON(http.StatusOK, getAllJobs())

	})

	r.GET("/job-history/:ID", func(c *gin.Context) {
		ID := c.Param("ID")
		job, success := getJobbyID(ID)

		if !success {
			c.JSON(http.StatusNotFound, gin.H{"message": "Job not found!"})
		} else {
			c.JSON(http.StatusOK, job)
		}
	})

	r.POST("/start-job", func(c *gin.Context) {
		var incomingJob coffeeJob
		c.BindJSON(&incomingJob)

		result, success, systemStatus := newJob(incomingJob.Product)

		if !success {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
			}).Debug("Error starting job")
			c.JSON(http.StatusBadRequest, gin.H{"message": systemStatus})
		} else {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
				"product":   result.Product,
				"jobReady":  result.JobReady,
			}).Debug("Job started")
			c.JSON(http.StatusOK, result)
		}

	})

	r.GET("/retrieve-job/:ID", func(c *gin.Context) {
		ID := c.Param("ID")
		job, success := retrieveJob(ID)
		log.WithFields(logrus.Fields{
			"requestId": c.MustGet("RequestId"),
			"success":   success,
		}).Debug("Request to retrieve Coffee")
		if !success {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
				"success":   success,
				"job-id":    ID,
			}).Debug("Failed to retrieve Coffee")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Nope!"})
		} else {
			c.JSON(http.StatusOK, job)
		}
	})

	return r
}

func main() {

	r := setupRouter()
	r.Run(":1337")

}
