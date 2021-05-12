package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/likakuli/generic-project-template/pkg/interfaces/mocks"
	"github.com/likakuli/generic-project-template/pkg/viewmodels"
)

func TestGetPlayerScore(t *testing.T) {
	playerService := mocks.IPlayerService{}

	playerService.On("GetScores", context.Background(), "Lucy", "Lily").Return("Forty-Fifteen", nil)
	playerController := NewPlayerController(&playerService)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/score/Lucy/vs/Lily", nil)
	w := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/api/v1/score/:player1/vs/:player2", playerController.GetPlayerScore)
	router.ServeHTTP(w, req)

	expectedResult := viewmodels.ScoreVM{}
	expectedResult.Score = "Forty-Fifteen"
	actualResult := viewmodels.ScoreVM{}

	json.NewDecoder(w.Body).Decode(&actualResult)

	// assert that the expectations were met
	assert.Equal(t, expectedResult, actualResult)
}
