package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/valyala/fastjson"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ApiResponse := events.APIGatewayProxyResponse{}
	// headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "*", "Access-Control-Allow-Methods": "*"}
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Headers"] = "*"
	headers["Access-Control-Allow-Methods"] = "*"
	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Credentials"] = "true"
	headers["Set-Cookie"] = "SameSite=Strict"

	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "POST":
		//validates json and returns error if not working
		err := fastjson.Validate(request.Body)
		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
		} else {
			_, err = dbHandler(request.Body)
			if err != nil {
				body := "Error: Unable to write to database ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
				ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
			}
			ApiResponse = events.APIGatewayProxyResponse{Headers: headers, Body: "Noice Message Sent", StatusCode: 200}
		}
	default:
		ApiResponse = events.APIGatewayProxyResponse{Body: "Method Not Allowed", StatusCode: 405}

	}
	// Response
	return ApiResponse, nil
}

type api_response struct {
	Name    string
	Email   string
	Message string
}

func dbHandler(body string) (*dynamodb.PutItemOutput, error) {

	var resp api_response
	fmt.Println(body)
	b := []byte(body)
	err := json.Unmarshal(b, &resp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	svc := dynamodb.New(session.New())
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Time": {
				S: aws.String(time.Now().Format("2006-01-02 15:04:05")),
			},
			"Name": {
				S: &resp.Name,
			},
			"Email": {
				S: &resp.Email,
			},
			"Message": {
				S: &resp.Message,
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	result, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				fmt.Println(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

	}
	return result, err

}

func main() {
	lambda.Start(HandleRequest)
}
