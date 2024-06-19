package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsClient *sqs.SQS
	queueURL  = os.Getenv("PAYMENTS_QUEUE_URL")
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

func init() {
	sess := session.Must(session.NewSession())
	sqsClient = sqs.New(sess)
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

	eventBody, err := json.Marshal(myEvent)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error marshalling event: %v", err)}, nil
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(eventBody)),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error sending SQS message: %v", err)}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(eventBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
