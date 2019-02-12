package main

import (
  "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestRootRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Hello World!", w.Body.String())
}

func TestHealthcheck(t *testing.T) {
  responsejson := gin.H{
    "message": "Ok",
  }

  router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
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

}

