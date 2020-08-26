package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	apirequest := &events.APIGatewayProxyRequest{}

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "hello", response.Body)
}
