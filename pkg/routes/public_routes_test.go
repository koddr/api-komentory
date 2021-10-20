package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPublicRoutes(t *testing.T) {
	// Load .env.test file from the root folder
	if err := godotenv.Load("../../.env.test"); err != nil {
		panic(err)
	}

	// Define a structure for specifying input and output data of a single test case.
	tests := []struct {
		description  string
		httpMethod   string
		route        string // input route
		expectedCode int
	}{
		// Successful test cases:
		{
			"success: get all projects",
			"GET", "/v1/projects",
			200,
		},
		{
			"fail: get all projects by not found user id",
			"GET", fmt.Sprintf("/v1/user/%s/projects", uuid.New().String()),
			200,
		},
		// Failed test cases:
		{
			"fail: get project by not found id",
			"GET", fmt.Sprintf("/v1/project/%s", uuid.New().String()),
			404,
		},
	}

	// Define Fiber app.
	app := fiber.New()

	// Define routes.
	PublicRoutes(app)

	// Iterate through test single test cases.
	for index, test := range tests {
		// Create a new http request with the route from the test case.
		req := httptest.NewRequest(test.httpMethod, test.route, nil)
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
		var resultMsg string
		if result["msg"] == nil {
			resultMsg = "no error message"
		} else {
			resultMsg = result["msg"].(string)
		}
		status := int(result["status"].(float64))
		description := fmt.Sprintf(
			"[%d] need to %s\nreal error output: %s",
			readableIndex, test.description, resultMsg,
		)

		// Checking, if the JSON field "status" from the response body has the expected status code.
		assert.Equalf(t, test.expectedCode, status, description)
	}
}
