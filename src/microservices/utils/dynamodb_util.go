package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dynamoDBClient *dynamodb.DynamoDB

func InitDynamoDBClient() {
	if dynamoDBClient == nil {
		sess := session.Must(session.NewSession())
		dynamoDBClient = dynamodb.New(sess)
	}
}

func SaveToDynamoDB(tableName, orderID string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	if dynamoDBClient == nil {
		InitDynamoDBClient()
	}

	av["order_id"] = &dynamodb.AttributeValue{
		S: aws.String(orderID),
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDBClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
