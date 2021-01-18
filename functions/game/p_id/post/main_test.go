package main

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "test")
	setup.AddOption(user, "test2")

	game, startResult := setup.Start(user)
	assert.Equal(t, theilliminationgame.Success, startResult)

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

func TestCancelRunningGame(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "test")
	setup.AddOption(user, "test2")

	game, startResult := setup.Start(user)
	assert.Equal(t, theilliminationgame.Success, startResult)

	id := game.Summary(user).ID

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = id.Hex()

	setupRequest := &SetUpRequest{
		UpdateType: "cancel",
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

	assert.Equal(t, id, gameResponse.Game.ID)
	assert.Equal(t, 2, len(gameResponse.Game.Remaining))

	assert.Equal(t, string(models.StateCancelled), gameResponse.Game.Status)
}

func TestReviewCompletedGame(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "test")
	setup.AddOption(user, "test2")

	game, startResult := setup.Start(user)
	assert.Equal(t, theilliminationgame.Success, startResult)
	game.Illiminate(user, "test")

	summary := game.Summary(user)
	assert.NotNil(t, summary.Winner)

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = summary.ID.Hex()

	setupRequest := `{
		"updateType":     "review",
		"reviewThoughts": "rubbish"
	}`

	apirequest.Body = string(setupRequest)

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var gameResponse GameResponse
	err = json.Unmarshal([]byte(response.Body), &gameResponse)
	assert.Nil(t, err)

	assert.Equal(t, summary.ID, gameResponse.Game.ID)

	assert.Equal(t, string(models.StateFinished), gameResponse.Game.Status)
	assert.NotNil(t, gameResponse.Game.CompletedGame.PlayerReviews)
	assert.Equal(t, "rubbish", gameResponse.Game.CompletedGame.PlayerReviews[0].Thoughts)
}

func TestSequentialReviewUpdatesDontOverwrite(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	user := illiminationtesting.TestUser(t, "User")
	setup := theilliminationgame.Create(user)
	setup.AddOption(user, "test")
	setup.AddOption(user, "test2")

	game, startResult := setup.Start(user)
	assert.Equal(t, theilliminationgame.Success, startResult)
	game.Illiminate(user, "test")

	summary := game.Summary(user)
	assert.NotNil(t, summary.Winner)

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("Test_User")
	apirequest.PathParameters = make(map[string]string)
	apirequest.PathParameters["id"] = summary.ID.Hex()

	setupRequest1 := `{
		"updateType":     "review",
		"reviewThoughts": "rubbish"
	}`
	setupRequest2 := `{
		"updateType":     "review",
		"hasImage": true
	}`

	// Act

	apirequest.Body = string(setupRequest1)
	response, err := Handler(*apirequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	apirequest.Body = string(setupRequest2)
	response, err = Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var gameResponse GameResponse
	err = json.Unmarshal([]byte(response.Body), &gameResponse)
	assert.Nil(t, err)

	assert.Equal(t, summary.ID, gameResponse.Game.ID)

	assert.Equal(t, string(models.StateFinished), gameResponse.Game.Status)
	assert.NotNil(t, gameResponse.Game.CompletedGame.PlayerReviews)
	assert.Equal(t, "rubbish", gameResponse.Game.CompletedGame.PlayerReviews[0].Thoughts)
}
