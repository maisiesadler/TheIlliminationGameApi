package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	logAttributes(&request)

	resp := ResponseSuccessful("Hello, world")
	return resp, nil
}

func logAttributes(request *events.APIGatewayProxyRequest) {

	for k, v := range request.RequestContext.Authorizer {
		fmt.Printf("Got k='%v', v='%v'\n", k, v)
	}

	if claims, ok := request.RequestContext.Authorizer["claims"]; ok {
		c := claims.(map[string]interface{})
		// username, ok := c["cognito:username"]
		// return username.(string), ok

		for k, v := range c {
			fmt.Printf("Claims has k='%v', v='%v'\n", k, v)
		}
	}

}

// ResponseSuccessful returns a 200 response for API Gateway that allows cors
func ResponseSuccessful(body string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,Authorization"
	resp.Body = body
	resp.StatusCode = 200
	return resp
}

func main() {
	lambda.Start(Handler)
}
