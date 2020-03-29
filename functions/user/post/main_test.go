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
	assert.Nil(t, err)
	assert.Equal(t, "Hello, world", response.Body)
}
