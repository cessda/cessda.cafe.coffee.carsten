// Carsten's CESSDA CAFE Coffee Machine
// Copyright CESSDA-ERIC 2022
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// CESSDA CAFE coffee machine implementation
package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/atarantini/ginrequestid"
	formatters "github.com/fabienm/go-logrus-formatters"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

		log.Debug(c.Request)

		log.WithFields(logrus.Fields{
			"method":    c.Request.Method,
			"requestId": c.MustGet("RequestId"),
			"endpoint":  path,
		}).Debug("request received")

		c.Next()

		elapsed := time.Since(start)

		log.WithFields(logrus.Fields{
			"method":       c.Request.Method,
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

	_, exists := os.LookupEnv("GELF_LOGGING")
	if exists {
		hostname, _ := os.Hostname()
		log.SetFormatter(formatters.NewGelf(hostname))
	}

	if gin.Mode() == "debug" {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}

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

		_, systemHTTPStatusCode, systemStatusMessage := systemStatus()

		c.JSON(systemHTTPStatusCode, gin.H{"message": systemStatusMessage})

	})

	r.GET("/metrics", func(c *gin.Context) {

		systemStatusCode, _, _ := systemStatus()

		c.String(http.StatusOK, "# HELP coffee_machine_status The system status 0=idle,1=brewing,2=blocked or other.\n# TYPE coffee_machine_status gauge\ncoffee_machine_status "+strconv.Itoa(systemStatusCode))

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

		result, success, systemHTTPStatusCode, systemStatusMessage := newJob(incomingJob.ID, incomingJob.Product)

		if !success {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
				"errorCode": systemHTTPStatusCode,
			}).Error("Error starting job: " + systemStatusMessage)
			c.JSON(systemHTTPStatusCode, gin.H{"message": systemStatusMessage})
		} else {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
				"product":   result.Product,
				"jobReady":  result.JobReady,
			}).Info("Job started")
			c.JSON(http.StatusOK, result)
		}

	})

	r.GET("/retrieve-job/:ID", func(c *gin.Context) {
		ID := c.Param("ID")
		job, success, systemHTTPStatusCode, message := retrieveJob(ID)
		log.WithFields(logrus.Fields{
			"requestId": c.MustGet("RequestId"),
			"success":   success,
		}).Info("Request to retrieve Coffee")
		if !success {
			log.WithFields(logrus.Fields{
				"requestId": c.MustGet("RequestId"),
				"errorCode": systemHTTPStatusCode,
				"job-id":    ID,
			}).Error("Failed to retrieve Coffee: " + message)
			c.JSON(http.StatusBadRequest, gin.H{"message": message})
		} else {
			c.JSON(http.StatusOK, job)
		}
	})

	// swagger:route POST /reset-button resetReq
	// Reset the Machine.
	// responses:
	//  200: succResp
	r.POST("/reset-button", func(c *gin.Context) {
		resetButton()
		log.WithFields(logrus.Fields{
			"requestId": c.MustGet("RequestId"),
		}).Error("Machine reset!")
		c.JSON(http.StatusOK, gin.H{"message": "Machine reset!"})
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
