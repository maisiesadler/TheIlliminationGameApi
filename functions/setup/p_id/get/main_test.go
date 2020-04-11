package main

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	id := setup.Summary(user).ID

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = id.Hex()

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var gamesResponse GamesResponse
	err = json.Unmarshal([]byte(response.Body), &gamesResponse)
	assert.Nil(t, err)

	assert.Equal(t, id, gamesResponse.Game.ID)
}
