package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/maisiesadler/theilliminationgame"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/maisiesadler/theilliminationgame/apigateway"
)

var errAuth = errors.New("Not logged in")
var errParse = errors.New("Error parsing request")
var errInvalidParameter = errors.New("Invalid parameter")

// StartGameRequest is the request from this handler
type StartGameRequest struct {
	SetUpID string `json:"setupId"`
}

// StartGameResponse is the response from this handler
type StartGameResponse struct {
	Result string                           `json:"result"`
	Game   *theilliminationgame.GameSummary `json:"game"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	user, err := apigateway.GetOrCreateAuthenticatedUser(context.TODO(), &request)
	if err != nil {
		return apigateway.ResponseUnsuccessful(401), errAuth
	}

	r := &StartGameRequest{}
	err = json.Unmarshal([]byte(request.Body), r)
	if err != nil {
		fmt.Printf("Could not parse body: %v.\n", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: 500}, errParse
	}

	objID, err := primitive.ObjectIDFromHex(r.SetUpID)
	if err != nil {
		return apigateway.ResponseUnsuccessful(400), errInvalidParameter
	}

	setup, err := theilliminationgame.LoadGameSetUp(&objID)
	if err != nil {
		fmt.Printf("error loading game setup: %v", err)
		return apigateway.ResponseUnsuccessful(500), err
	}

	game, startResult := setup.Start(user)
	var summary *theilliminationgame.GameSummary
	if startResult == theilliminationgame.Success {
		summary = game.Summary(user)
	}

	response := &StartGameResponse{
		Result: string(startResult),
		Game:   summary,
	}

	resp := apigateway.ResponseSuccessful(response)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
