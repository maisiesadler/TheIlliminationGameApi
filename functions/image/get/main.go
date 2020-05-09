package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/reviewphotos"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errMissingParameters = errors.New("Missing parameters")
var errParseResult = errors.New("Error parsing result")

// ImageURLResponse is the type returned by get image url
type ImageURLResponse struct {
	URL string `json:"url"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request SearchIngredient %s\n", request.RequestContext.RequestID)

	verb, ok := request.QueryStringParameters["verb"]
	if !ok {
		fmt.Println("Missing parameter 'verb'.")
		return apigateway.ResponseUnsuccessful(400), errMissingParameters
	}
	key, ok := request.QueryStringParameters["key"]
	if !ok {
		fmt.Println("Missing parameter 'key'.")
		return apigateway.ResponseUnsuccessful(400), errMissingParameters
	}
	if _, appenduser := request.QueryStringParameters["appenduser"]; appenduser {
		user, err := apigateway.GetOrCreateAuthenticatedUser(context.TODO(), &request)
		if err != nil {
			return apigateway.ResponseUnsuccessful(401), nil
		}

		// append userID so each user can have their unique photo
		key = key + "_" + user.ViewID.Hex()
	}

	url, err := reviewphotos.CreatePresignedURL(verb, key)
	if err != nil {
		log.Fatal(err.Error())
	}

	response := &ImageURLResponse{
		URL: url,
	}

	return apigateway.ResponseSuccessful(response), nil
}

func main() {
	lambda.Start(Handler)
}
