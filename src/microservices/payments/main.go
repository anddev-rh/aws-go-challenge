package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"microservices/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsClient         *sqs.SQS
	queueURL          = os.Getenv("ORDERS_QUEUE_URL")
	paymentsTableName = os.Getenv("PAYMENTS_TABLE_NAME")
)

type ProcessPaymentsRequest struct {
	OrderID string       `json:"order_id"`
	Status  utils.Status `json:"status"`
}

func init() {
	sess := session.Must(session.NewSession())
	sqsClient = sqs.New(sess)
}

func validateProcessPaymentsRequest(body []byte) (*ProcessPaymentsRequest, error) {
	var processPaymentsRequest ProcessPaymentsRequest
	err := json.Unmarshal(body, &processPaymentsRequest)
	if err != nil {
		return nil, err
	}

	if processPaymentsRequest.OrderID == "" {
		return nil, errors.New("invalid order_id")
	}
	if processPaymentsRequest.Status != utils.StatusIncompleted {
		return nil, errors.New("invalid status")
	}
	return &processPaymentsRequest, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	validatedRequest, err := validateProcessPaymentsRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	validatedRequest.Status = utils.StatusCompleted

	err = utils.SaveToDynamoDB(paymentsTableName, validatedRequest.OrderID, validatedRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error saving payment to DynamoDB: %v", err),
		}, nil
	}

	eventBody, err := json.Marshal(validatedRequest)
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
