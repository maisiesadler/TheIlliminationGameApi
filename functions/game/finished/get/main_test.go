package main

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameFindPredicate(func(setup *models.Game, m primitive.M) bool { return true })

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var gamesResponse GamesResponse
	err = json.Unmarshal([]byte(response.Body), &gamesResponse)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(gamesResponse.Games))
}
