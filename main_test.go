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

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

// Validate the response for root path
func TestRootRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome to the CESSDA Café!", w.Body.String())
}

// Validate HealthCeck path
func TestHealthcheck(t *testing.T) {
	var requestUUID string
	requestUUID = "11122333-aaaa-bbbb-cccc-444555666777"

	responsejson := gin.H{
		"message": "Ok",
	}

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	req.Header.Add("X-Request-Id", requestUUID)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	// Grab the value & whether or not it exists
	value, exists := response["message"]

	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, responsejson["message"], value)

	// check request header
	assert.Contains(t, req.Header["X-Request-Id"], requestUUID)

}
