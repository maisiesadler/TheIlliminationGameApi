package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/maisiesadler/theilliminationgame/apigateway"
)

var errAuth = errors.New("Not logged in")
var errParse = errors.New("Error parsing response")
var errMissingParameter = errors.New("Missing parameter")
var errInvalidParameter = errors.New("Invalid parameter")

// GamesResponse is the response from this handler
type GamesResponse struct {
	Game *theilliminationgame.GameSetUpSummary `json:"game"`
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

	id, ok := request.PathParameters["id"]
	if !ok || id == "" {
		return apigateway.ResponseUnsuccessful(400), errMissingParameter
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apigateway.ResponseUnsuccessful(400), errInvalidParameter
	}

	game, err := theilliminationgame.LoadGameSetUp(&objID)
	if err != nil {
		fmt.Printf("Error finding games: '%v'.\n", err)
		return apigateway.ResponseUnsuccessful(500), err
	}

	response := &GamesResponse{
		Game: game.Summary(user),
	}

	resp := apigateway.ResponseSuccessful(response)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
