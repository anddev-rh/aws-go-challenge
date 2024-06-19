package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	testCases := []struct {
		name          string
		request       events.APIGatewayProxyRequest
		expectedBody  events.APIGatewayProxyResponse
		expectedError error
	}{
		{
			name: "Valid request",
			request: events.APIGatewayProxyRequest{
				Body: `{order_id: "123", status: "pending"}`,
			},
			expectedBody:  events.APIGatewayProxyResponse{StatusCode: 200, Body: `{"order_id":"123","total_price":100}`},
			expectedError: nil,
		},
		{
			name: "Invalid request",
			request: events.APIGatewayProxyRequest{
				Body: `{"user_id": "", "item": "item1", "quantity": 1, "total_price": 100}`,
			},
			expectedBody:  events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body: invalid user_id"},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := handler(testCase.request)

			if err != testCase.expectedError {
				t.Errorf("Expected error %v, but got %v", testCase.expectedError, err)
			}
		})
	}
}
