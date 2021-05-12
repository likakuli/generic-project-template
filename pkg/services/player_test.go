package services

import (
	"context"
	"testing"

	"github.com/likakuli/generic-project-template/pkg/interfaces/mocks"
	"github.com/likakuli/generic-project-template/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestGetScores(t *testing.T) {
	playerRepo := mocks.IPlayerRepository{}

	player1 := models.Player{}
	player1.ID = 101
	player1.Name = "Lucy"
	player1.Score = 3

	player2 := models.Player{}
	player2.ID = 102
	player2.Name = "Lily"
	player2.Score = 1

	playerRepo.On("GetPlayerByName", "Lucy").Return(&player1, nil)
	playerRepo.On("GetPlayerByName", "Lily").Return(&player2, nil)

	playerService := NewPlayerService(&playerRepo)

	expectedResult := "Forty-Fifteen"
	actualResult, _ := playerService.GetScores(context.Background(), "Lucy", "Lily")
	assert.Equal(t, expectedResult, actualResult)
}
