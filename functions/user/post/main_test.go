package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maisiesadler/theilliminationgame/illiminationtesting"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange

	illiminationtesting.SetTestCollectionOverride()

	apirequest := &events.APIGatewayProxyRequest{}

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, "", response.Body)
}
