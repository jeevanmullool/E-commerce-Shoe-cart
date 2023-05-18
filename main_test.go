package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"redkart/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminLogin(t *testing.T) {
	// Create a new instance of the Gin engine
	router := setupRouter()

	// Define the JSON payload for the request
	requestBody := map[string]string{
		"Email":    "admin@example.com",
		"Password": "password123",
	}
	jsonPayload, _ := json.Marshal(requestBody)

	// Create a new HTTP POST request with the JSON payload
	req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rec := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rec, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert that the response body is empty ({} JSON object)
	assert.Equal(t, "{}", rec.Body.String())
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/admin/login", controllers.AdminLogin)

	return router
}
