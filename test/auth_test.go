package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jabutech/simple-blog/auth"
	"github.com/jabutech/simple-blog/router"
	"github.com/jabutech/simple-blog/util"
	"github.com/stretchr/testify/assert"
)

// Func for create account random
func createRandomAccount(t *testing.T, withIsAdmin bool) auth.RegisterInput {
	// Open connection to db
	db := util.SetupTestDb()

	// Call router with argument db
	router := router.SetupRouter(db)

	var data auth.RegisterInput
	var dataBody string

	if !withIsAdmin {
		data = auth.RegisterInput{
			Fullname: util.RandomFullname(),
			Email:    util.RandomEmail(),
			Password: "password",
		}

		dataBody = fmt.Sprintf(`{"fullname": "%s", "email": "%s", "password": "%s"}`, data.Fullname, data.Email, "password")
	} else {

		data = auth.RegisterInput{
			Fullname: util.RandomFullname(),
			Email:    util.RandomEmail(),
			Password: "password",
			IsAdmin:  true,
		}

		dataBody = fmt.Sprintf(`{"fullname": "%s", "email": "%s", "password": "%s", "is_admin": %t}`, data.Fullname, data.Email, "password", true)
	}

	// Create payload request
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/register", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run server http
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()

	// Read response
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	// Decode json
	json.Unmarshal(body, &responseBody)

	// Response status code must be 200 (success)
	assert.Equal(t, 200, response.StatusCode)
	// Response body status code must be 200 (success)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// Response body status must be success
	assert.Equal(t, "success", responseBody["status"])
	// Response body message
	assert.Equal(t, "You have successfully registered.", responseBody["message"])

	// Return data success register
	return data
}

// Test Register
func TestRegisterSuccessWithoutIsAdmin(t *testing.T) {
	// Var withIsAdmin value false
	withIsAdmin := false
	// Call function createRandomAccount for test create account
	createRandomAccount(t, withIsAdmin)
}

func TestRegisterSuccessWithIsAdmin(t *testing.T) {

	// Var withIsAdmin value true
	withIsAdmin := true
	// Call function createRandomAccount for test create account
	createRandomAccount(t, withIsAdmin)

}

func TestRegisterValidationError(t *testing.T) {
	// Open connection to db
	db := util.SetupTestDb()

	// Call router with argument db
	router := router.SetupRouter(db)

	dataBody := fmt.Sprintf(`{"fullname": "", "email": "aa", "password": "a" }`)

	// Create payload request
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/register", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run server http
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()

	// Read response
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	// Decode json
	json.Unmarshal(body, &responseBody)

	// Response status code must be 200 (success)
	assert.Equal(t, 400, response.StatusCode)
	// Response body status code must be 200 (success)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	// Response body status must be success
	assert.Equal(t, "error", responseBody["status"])
	// Response body message
	assert.Equal(t, "Registered failed.", responseBody["message"])
}

// loginRandomAccount for login simulation with random account that returns jwt token
func loginRandomAccount(t *testing.T) interface{} {
	// Open connection to db
	db := util.SetupTestDb()

	// Var withIsAdmin value false
	withIsAdmin := false

	// Create account random
	account := createRandomAccount(t, withIsAdmin)
	// Call router with argument db
	router := router.SetupRouter(db)

	// Data body with data from create account random
	dataBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, account.Email, account.Password)

	// Create payload request
	requestBody := strings.NewReader(dataBody)

	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/login", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run server http
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()

	// Read response
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	// Decode json
	json.Unmarshal(body, &responseBody)

	// Response status code must be 200 (success)
	assert.Equal(t, 200, response.StatusCode)
	// Response body status code must be 200 (success)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// Response body status must be success
	assert.Equal(t, "success", responseBody["status"])
	// Response body message
	assert.Equal(t, "You have Login.", responseBody["message"])
	// Data is not null
	assert.NotZero(t, responseBody["data"])
	// Property token is not nul
	assert.NotZero(t, responseBody["data"].(map[string]interface{})["token"])

	// Get token for return use any test
	token := responseBody["data"].(map[string]interface{})["token"]
	return token
}

// Test Login
func TestLoginSuccess(t *testing.T) {
	loginRandomAccount(t)
}

// Test get all users
func TestGetListUsers(t *testing.T) {
	// Open connection to db
	db := util.SetupTestDb()

	// Call router with argument db
	router := router.SetupRouter(db)

	// login for get token
	token := loginRandomAccount(t)
	strToken := fmt.Sprintf("Bearer %s", token)

	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/users", nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization with by inserting jwt token
	request.Header.Add("Authorization", strToken)

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run server http
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()

	// Read response
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	// Decode json
	json.Unmarshal(body, &responseBody)
	// Response status code must be 200 (success)
	assert.Equal(t, 200, response.StatusCode)
	// Response body status code must be 200 (success)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// Response body status must be success
	assert.Equal(t, "success", responseBody["status"])
	// Response body message
	assert.Equal(t, "List of users", responseBody["message"])

	// Response data list users
	var listUsers = responseBody["data"].([]interface{})
	// Response body data length not 0
	assert.NotEqual(t, 0, len(listUsers))

	// All property not empty
	for _, list := range listUsers {
		mapList := list.(map[string]interface{})
		assert.NotEmpty(t, mapList["id"])
		assert.NotEmpty(t, mapList["email"])
		assert.NotEmpty(t, mapList["fullname"])
		assert.NotNil(t, mapList["is_admin"])
	}
}
