// Carsten's CESSDA CAFE Coffee Machine
// Copyright CESSDA-ERIC 2019
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// CESSDA CAFE coffee machine implementation
package main

import (
	"github.com/atarantini/ginrequestid"
	"github.com/gin-gonic/gin"
	//gelf "github.com/seatgeek/logrus-gelf-formatter"
	"github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type coffeeJob struct {
	ID      string `json:"jobId"`
	Product string `json:"product"`
}

// Success response
// swagger:response succResp
type swaggSuccResp struct {
	// in:body
	Body struct {
		// Detailed message
		Message string `json:"message"`
	}
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
	//log.Formatter = new(gelf.GelfFormatter)

	graylogServer, exists := os.LookupEnv("GRAYLOG_SERVER")
	if exists {
		hook := graylog.NewGraylogHook(graylogServer, map[string]interface{}{"source": "coffee-carsten"})
		log.AddHook(hook)
	}

	log.Level = logrus.DebugLevel

	r.Use(myRequestLogger(log), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the CESSDA Café!")
	})

	// swagger:route GET /healthcheck healthReq
	// Perform a healthcheck.
	// If all is fine, this will tell you.
	// responses:
	//  200: succResp
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

		result, success, systemStatus := newJob(incomingJob.ID, incomingJob.Product)

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

// Get environment variable with a default
func getEnv(name string, defaultValue string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	}
	return defaultValue
}

func main() {

	coffeePort := getEnv("COFFEE_PORT", "1337")

	r := setupRouter()
	r.Run(":" + coffeePort)

}
