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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsClient         *sqs.SQS
	queueURL          = os.Getenv("PAYMENTS_QUEUE_URL")
	dynamoDBClient    *dynamodb.DynamoDB
	paymentsTableName = os.Getenv("PAYMENTS_TABLE_NAME")
)

type Status string

const (
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusPending   Status = "pending"
)

type ProcessPaymentsRequest struct {
	OrderID string `json:"order_id"`
	Status  Status `json:"status"`
}

func init() {
	sess := session.Must(session.NewSession())
	sqsClient = sqs.New(sess)
	dynamoDBClient = dynamodb.New(sess)
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
	if processPaymentsRequest.Status != StatusPending {
		return nil, errors.New("invalid status")
	}
	return &processPaymentsRequest, nil
}

func savePaymentToDynamoDB(processPaymentsRequest *ProcessPaymentsRequest) error {
	av, err := dynamodbattribute.MarshalMap(processPaymentsRequest)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(paymentsTableName),
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	validatedRequest, err := validateProcessPaymentsRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	validatedRequest.Status = StatusCompleted

	err = savePaymentToDynamoDB(validatedRequest)
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
