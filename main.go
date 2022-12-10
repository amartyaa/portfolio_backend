package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/valyala/fastjson"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ApiResponse := events.APIGatewayProxyResponse{}
	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "POST":
		//validates json and returns error if not working
		err := fastjson.Validate(request.Body)
		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
		} else {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Message Sent", StatusCode: 200}
		}
	default:
		ApiResponse = events.APIGatewayProxyResponse{Body: "Method Not Allowed", StatusCode: 405}

	}
	// Response
	return ApiResponse, nil
}

func dbHandler() {
	
}

func main() {
	lambda.Start(HandleRequest)
}
