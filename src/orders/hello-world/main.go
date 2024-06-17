package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
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
	fmt.Println("err:", err, "validatedRequest:", validatedRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %v", err),
		}, nil
	}

	myEvent := CreateOrderEvent{
		OrderID:    uuid.New().String(),
		TotalPrice: validatedRequest.TotalPrice,
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
