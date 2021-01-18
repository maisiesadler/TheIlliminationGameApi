package main

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetUserViewFindPredicate(func(uv *models.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})
	illiminationtesting.SetUserOptionsFindPredicate(func(uo *models.UserOption, m primitive.M) bool {
		return m["userid"] == uo.UserID
	})

	apirequest := illiminationtesting.CreateTestAuthorizedRequest("User")

	// Act
	response, err := Handler(*apirequest)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	var uoResponse UserOptionsResponse
	err = json.Unmarshal([]byte(response.Body), &uoResponse)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(uoResponse.Options))
}
