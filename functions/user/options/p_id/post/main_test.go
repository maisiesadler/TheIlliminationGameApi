package main

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = primitive.NewObjectID().String()

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, 400, response.StatusCode)
}
