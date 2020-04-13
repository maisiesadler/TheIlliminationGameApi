package main

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetUserViewFindPredicate(func(uv *models.UserView, m bson.M) bool {
		return m["username"] == uv.Username
	})

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "One")
	setup.AddOption(user, "Two")
	setup.AddOption(user, "Three")

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	startRequest := &StartGameRequest{
		SetUpID: setup.Summary(user).ID.Hex(),
	}

	b, err := json.Marshal(startRequest)
	assert.Nil(t, err)
	apirequest.Body = string(b)

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var gameResponse StartGameResponse
	err = json.Unmarshal([]byte(response.Body), &gameResponse)
	assert.Nil(t, err)

	assert.Equal(t, string(theilliminationgame.Success), gameResponse.Result)

	if gameResponse.Game == nil {
		t.Errorf("Game is nil")
		t.FailNow()
	}

	assert.NotNil(t, gameResponse.Game.ID)
	assert.NotNil(t, gameResponse.Game.Remaining)
	assert.NotNil(t, gameResponse.Game.Illiminated)
	assert.NotNil(t, gameResponse.Game.Players)

	assert.Equal(t, 3, len(gameResponse.Game.Remaining))
	assert.Equal(t, 0, len(gameResponse.Game.Illiminated))
	assert.Equal(t, 1, len(gameResponse.Game.Players))
}
