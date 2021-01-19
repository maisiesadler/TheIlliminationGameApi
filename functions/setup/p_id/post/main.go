package main

import (
	"context"
	"encoding/json"
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
var errParse = errors.New("Error parsing request")
var errMissingParameter = errors.New("Missing parameter")
var errInvalidParameter = errors.New("Invalid parameter")

// SetUpRequest is the request for this handler
type SetUpRequest struct {
	UpdateType  string            `json:"updateType"`
	Option      string            `json:"option"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	OptionIndex int               `json:"optionIdx"`
	Updates     map[string]string `json:"updates"`
	Tag         *string           `json:"tag"`
}

// GameUpdateResponse is the response from this handler
type GameUpdateResponse struct {
	Game   *theilliminationgame.GameSetUpSummary `json:"game"`
	Result string                                `json:"result"`
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

	r := &SetUpRequest{}
	err = json.Unmarshal([]byte(request.Body), r)
	if err != nil {
		fmt.Printf("Could not parse body: %v.\n", request.Body)
		return events.APIGatewayProxyResponse{StatusCode: 500}, errParse
	}

	id, ok := request.PathParameters["id"]
	if !ok || id == "" {
		return apigateway.ResponseUnsuccessful(400), errMissingParameter
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apigateway.ResponseUnsuccessful(400), errInvalidParameter
	}

	setup, err := theilliminationgame.LoadGameSetUp(&objID)
	if err != nil {
		fmt.Printf("Error finding games: '%v'.\n", err)
		return apigateway.ResponseUnsuccessful(500), err
	}

	var result string

	if r.UpdateType == "join" {
		setup.JoinGame(user)
	} else if r.UpdateType == "option_add" {
		result = string(setup.AddOption(user, r.Option))
	} else if r.UpdateType == "detailedoption_add" {
		result = string(setup.AddDetailedOption(user, r.Option, r.Description, r.Link))
	} else if r.UpdateType == "option_update" {
		if r.Updates == nil {
			return apigateway.ResponseUnsuccessfulString(400, "No Updates"), err
		}
		setup.UpdateOption(user, r.OptionIndex, r.Updates)
	} else if r.UpdateType == "deactivate" {
		setup.Deactivate(user)
	} else if r.UpdateType == "addtag" {
		setup.AddTag(user, *r.Tag)
	} else if r.UpdateType == "removetag" {
		setup.RemoveTag(user, *r.Tag)
	} else {
		result = "Unknown update type"
	}

	setup, _ = theilliminationgame.LoadGameSetUp(&objID)

	response := &GameUpdateResponse{
		Result: result,
		Game:   setup.Summary(user),
	}

	resp := apigateway.ResponseSuccessful(response)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
