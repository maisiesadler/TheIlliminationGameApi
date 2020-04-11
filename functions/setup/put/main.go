package main

import (
	"context"
	"errors"
	"log"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/maisiesadler/theilliminationgame/apigateway"
)

var errAuth = errors.New("Not logged in")
var errParse = errors.New("Error parsing request")

// CreateGameResponse is the response from this handler
type CreateGameResponse struct {
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

	setup := theilliminationgame.Create(user)

	response := &CreateGameResponse{
		Game: setup.Summary(user),
	}

	resp := apigateway.ResponseSuccessful(response)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
