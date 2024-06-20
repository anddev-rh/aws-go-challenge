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
	"github.com/google/uuid"
)

var (
	sqsClient        *sqs.SQS
	paymentsQueueUrl = os.Getenv("PAYMENTS_QUEUE_URL")
	ordersTableName  = os.Getenv("ORDERS_TABLE_NAME")
)

type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

type CreateOrderEvent struct {
	OrderID    string       `json:"order_id"`
	TotalPrice int64        `json:"total_price"`
	Status     utils.Status `json:"status"`
}

func init() {
	sess := session.Must(session.NewSession())
	sqsClient = sqs.New(sess)
}

func validateCreateOrderRequest(body []byte) (*CreateOrderRequest, error) {
	var createOrderRequest CreateOrderRequest
	err := json.Unmarshal(body, &createOrderRequest)
	if err != nil {
		return nil, err
	}

	if createOrderRequest.UserID == "" {
		return nil, errors.New("invalid user_id")
	}
	if createOrderRequest.Item == "" {
		return nil, errors.New("invalid item")
	}
	if createOrderRequest.Quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}
	if createOrderRequest.TotalPrice <= 0 {
		return nil, errors.New("invalid total_price")
	}

	return &createOrderRequest, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	validatedRequest, err := validateCreateOrderRequest([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	myEvent := CreateOrderEvent{
		OrderID:    uuid.New().String(),
		TotalPrice: validatedRequest.TotalPrice,
		Status:     utils.StatusIncompleted,
	}

	err = utils.SaveToDynamoDB(ordersTableName, myEvent.OrderID, validatedRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error saving order to DynamoDB: %v", err),
		}, nil
	}

	eventBody, err := json.Marshal(myEvent)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error marshalling event: %v", err)}, nil
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(paymentsQueueUrl),
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
