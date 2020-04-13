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
	id := setup.Summary(user).ID

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = id.Hex()

	setupRequest := &SetUpRequest{
		UpdateType: "option",
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
	var gamesResponse GameUpdateResponse
	err = json.Unmarshal([]byte(response.Body), &gamesResponse)
	assert.Nil(t, err)

	assert.Equal(t, id, gamesResponse.Game.ID)
	assert.Equal(t, 1, len(gamesResponse.Game.Options))
	assert.Equal(t, "test", gamesResponse.Game.Options[0])
}
