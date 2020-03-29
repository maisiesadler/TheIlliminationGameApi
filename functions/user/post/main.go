package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/maisiesadler/theilliminationgame/apigateway"
)

var errAuth = errors.New("Not logged in")
var errParse = errors.New("Error parsing request")

// UpdateRecipeRequest is the request from this handler
type UpdateRecipeRequest struct {
	Nickname string `json:"nickname"`
}

// User is the response from this handler
type User struct {
	Nickname string `json:"nickname"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	r := &UpdateRecipeRequest{}
	err := json.Unmarshal([]byte(request.Body), r)
	if err != nil {
		fmt.Printf("Could not parse body: %v.\n", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: 500}, errParse
	}

	user, err := apigateway.GetOrCreateAuthenticatedUser(context.TODO(), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 401}, errAuth
	}

	err = user.SetNickname(context.TODO(), r.Nickname)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, errAuth
	}

	user, err = apigateway.GetOrCreateAuthenticatedUser(context.TODO(), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	response := &User{
		Nickname: user.Nickname,
	}

	resp := apigateway.ResponseSuccessful(response)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
