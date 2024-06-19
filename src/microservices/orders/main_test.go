package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name          string
		request       events.APIGatewayProxyRequest
		expectedCode  int
		expectedBody  string
		expectedError error
	}{
		{
			name: "Valid request",
			request: events.APIGatewayProxyRequest{
				Body: `{"user_id": "123", "item": "item1", "quantity": 1, "total_price": 100}`,
			},
			expectedBody:  `{"order_id":"123","total_price":100}`,
			expectedCode:  200,
			expectedError: nil,
		},
		{
			name: "Invalid request - missing user_id",
			request: events.APIGatewayProxyRequest{
				Body: `{"user_id": "", "item": "item1", "quantity": 1, "total_price": 100}`,
			},
			expectedCode:  400,
			expectedBody:  `Invalid request body: invalid user_id`,
			expectedError: nil,
		},
		{
			name: "Invalid request - negative quantity",
			request: events.APIGatewayProxyRequest{
				Body: `{"user_id": "123", "item": "item1", "quantity": -1, "total_price": 100}`,
			},
			expectedCode:  400,
			expectedBody:  `Invalid request body: invalid quantity`,
			expectedError: nil,
		},
		{
			name: "Invalid request - zero total_price",
			request: events.APIGatewayProxyRequest{
				Body: `{"user_id": "123", "item": "item1", "quantity": 1, "total_price": 0}`,
			},
			expectedCode:  400,
			expectedBody:  `Invalid request body: invalid total_price`,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := handler(testCase.request)

			if err != testCase.expectedError {
				t.Errorf("Expected error %v, but got %v", testCase.expectedError, err)
			}

			if response.StatusCode != testCase.expectedCode {
				t.Errorf("Expected status code %d, but got %d", testCase.expectedCode, response.StatusCode)
			}

			// Check response body for valid request
			if testCase.name == "Valid request" {
				var responseBody map[string]interface{}
				fmt.Println(response)
				if err := json.Unmarshal([]byte(response.Body), &responseBody); err != nil {
					t.Fatalf("Error unmarshalling actual body: %v", err)
				}

				orderID, ok := responseBody["order_id"].(string)
				if !ok {
					t.Error("Expected order_id to be present, but it was missing or not a string")
				}

				if orderID == "" {
					t.Error("Expected non-empty order_id in response body")
				}

				totalPrice, ok := responseBody["total_price"].(float64)
				if !ok {
					t.Error("Expected total_price to be present, but it was missing or not a number")
				}

				if totalPrice != 100 {
					t.Errorf("Expected total_price to be 100, but got %v", totalPrice)
				}
			} else {
				// Check response body for other cases
				if response.Body != testCase.expectedBody {
					t.Errorf("Expected body %s, but got %s", testCase.expectedBody, response.Body)
				}
			}
		})
	}
}
