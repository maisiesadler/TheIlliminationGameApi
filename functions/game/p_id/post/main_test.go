package main

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetUserViewFindPredicate(func(uv *models.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "test")
	setup.AddOption(user, "test2")
	game, err := setup.Start(user)
	assert.Nil(t, err)

	id := game.Summary(user).ID

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = id.Hex()

	setupRequest := &SetUpRequest{
		UpdateType: "illiminate",
		Option:     "test",
	}

	b, err := json.Marshal(setupRequest)
	assert.Nil(t, err)
	apirequest.Body = string(b)

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var gameResponse GameResponse
	err = json.Unmarshal([]byte(response.Body), &gameResponse)
	assert.Nil(t, err)

	assert.Equal(t, "Illiminated", gameResponse.Result)

	assert.Equal(t, id, gameResponse.Game.ID)
	assert.Equal(t, 1, len(gameResponse.Game.Remaining))
	assert.Equal(t, "test2", gameResponse.Game.Remaining[0])
	assert.Equal(t, 1, len(gameResponse.Game.Illiminated))
	assert.Equal(t, "test", gameResponse.Game.Illiminated[0])

	assert.Equal(t, string(models.StateFinished), gameResponse.Game.Status)
}
