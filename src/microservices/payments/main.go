package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Status string

const (
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusPending   Status = "pending"
)

type ProcesspaymentsRequest struct {
	OrderID string `json:"order_id"`
	Status  Status `json:"status"`
}

func validateProcesspaymentsRequest(body []byte) (*ProcesspaymentsRequest, error) {
	var processpaymentsRequest ProcesspaymentsRequest
	err := json.Unmarshal(body, &processpaymentsRequest)
	if err != nil {
		return nil, err
	}

	if processpaymentsRequest.OrderID == "" {
		return nil, errors.New("invalid order_id")
	}
	if processpaymentsRequest.Status != StatusPending {
		return nil, errors.New("invalid status")
	}
	return &processpaymentsRequest, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	validatedRequest, err := validateProcesspaymentsRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	myEvent := ProcesspaymentsRequest{
		OrderID: validatedRequest.OrderID,
		Status:  StatusCompleted,
	}

	myEventBytes, _ := json.Marshal(myEvent)
	return events.APIGatewayProxyResponse{
		Body:       string(myEventBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
