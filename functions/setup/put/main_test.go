package main

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var gameResponse CreateGameResponse
	err = json.Unmarshal([]byte(response.Body), &gameResponse)
	assert.Nil(t, err)

	assert.NotNil(t, gameResponse.Game.Code)
	assert.NotNil(t, gameResponse.Game.ID)
	assert.NotNil(t, gameResponse.Game.Options)
	assert.NotNil(t, gameResponse.Game.Players)

	assert.Equal(t, 0, len(gameResponse.Game.Options))
	assert.Equal(t, 1, len(gameResponse.Game.Players))
}
