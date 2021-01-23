package main

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = "12345"

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, 400, response.StatusCode)
}
