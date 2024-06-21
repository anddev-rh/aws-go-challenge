package main

import (
	"encoding/json"
	"fmt"
	"os"

	"microservices/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	dynamoDBClient  *dynamodb.DynamoDB
	ordersTableName = os.Getenv("ORDERS_TABLE_NAME")
)

type ProcessPaymentsRequest struct {
	OrderID string       `json:"order_id"`
	Status  utils.Status `json:"status"`
}

func init() {
	sess := session.Must(session.NewSession())
	dynamoDBClient = dynamodb.New(sess)
}

func updateOrderStatus(orderID string, status string) error {
	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ordersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"order_id": {
				S: aws.String(orderID),
			},
		},
	})
	if err != nil {
		return err
	}

	item := make(map[string]*dynamodb.AttributeValue)
	for key, value := range result.Item {
		item[key] = value
	}

	item["status"] = &dynamodb.AttributeValue{
		S: aws.String(status),
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(ordersTableName),
		Item:      item,
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Updated order %s status to %s\n", orderID, status)
	return nil
}

func handler(request events.SQSEvent) error {
	for _, record := range request.Records {
		var event ProcessPaymentsRequest
		err := json.Unmarshal([]byte(record.Body), &event)
		if err != nil {
			fmt.Printf("Error unmarshalling SQS message body: %v\n", err)
			continue
		}

		if event.Status != utils.StatusCompleted {
			fmt.Printf("Received invalid status in SQS message: %s\n", event.Status)
			continue
		}

		err = updateOrderStatus(event.OrderID, string(utils.StatusReady))
		if err != nil {
			fmt.Printf("Error updating order status: %v\n", err)
			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
