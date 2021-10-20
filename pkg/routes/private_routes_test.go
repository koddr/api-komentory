package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPrivateRoutes(t *testing.T) {
	// Load .env.test file from the root folder.
	if err := godotenv.Load("../../.env.test"); err != nil {
		panic(err)
	}

	// Define test variables.
	body := map[string]string{
		"empty":     `{}`,
		"non-empty": `{"title": "Test title"}`,
	}

	// Define a structure for specifying input and output data of a single test case.
	tests := []struct {
		description  string
		method       string // input method
		route        string // input route
		tokenString  string // input token
		body         io.Reader
		expectedCode int
	}{
		// Failed test cases:
		{
			"fail: create project without JWT and JSON body",
			"POST", "/v1/project", "", bytes.NewBuffer([]byte(body["empty"])),
			400, // Missing or malformed JWT
		},
		{
			"fail: update project without JWT and JSON body",
			"PATCH", "/v1/project", "", bytes.NewBuffer([]byte(body["empty"])),
			400, // Missing or malformed JWT
		},
		{
			"fail: delete project without JWT and JSON body",
			"DELETE", "/v1/project", "", bytes.NewBuffer([]byte(body["empty"])),
			400, // Missing or malformed JWT
		},
		{
			"fail: put file to CDN without JWT and JSON body",
			"PUT", "/v1/cdn/upload", "", bytes.NewBuffer([]byte(body["empty"])),
			400, // Missing or malformed JWT
		},
	}

	// Define a new Fiber app.
	app := fiber.New()

	// Define routes.
	PrivateRoutes(app)

	// Iterate through test single test cases
	for index, test := range tests {
		// Create a new http request with the route from the test case.
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.tokenString))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request plain with the app.
		resp, _ := app.Test(req, -1) // the -1 disables request latency

		// Parse the response body.
		body, errReadAll := io.ReadAll(resp.Body)
		if errReadAll != nil {
			return
		}

		// Set the response body (JSON) to simple map.
		var result map[string]interface{}
		if errUnmarshal := json.Unmarshal(body, &result); errUnmarshal != nil {
			return
		}

		// Redefine index of the test case.
		readableIndex := index + 1

		// Define status & description from the response.
		status := int(result["status"].(float64))
		description := fmt.Sprintf(
			"[%d] need to %s\nreal error output: %s",
			readableIndex, test.description, result["msg"].(string),
		)

		// Checking, if the JSON field "status" from the response body has the expected status code.
		assert.Equalf(t, test.expectedCode, status, description)
	}
}
