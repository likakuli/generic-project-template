package services

import (
	"context"

	"github.com/leopoldxx/go-utils/trace"

	"github.com/likakuli/generic-project-template/pkg/interfaces"
)

type playerService struct {
	repo interfaces.IPlayerRepository
}

func (service *playerService) GetScores(ctx context.Context, player1Name string, player2Name string) (string, error) {
	tracer := trace.GetTraceFromContext(ctx)

	baseScore := [4]string{"Love", "Fifteen", "Thirty", "Forty"}
	var result string

	player1, err := service.repo.GetPlayerByName(player1Name)
	if err != nil {
		tracer.Infof("GetPlayerByName for play1: %s failed with err: %s", player1Name, err.Error())
		return "", err
	}

	player2, err := service.repo.GetPlayerByName(player2Name)
	if err != nil {
		tracer.Infof("GetPlayerByName for play2: %s failed with err: %s", player2Name, err.Error())
		return "", err
	}

	if player1.Score < 4 && player2.Score < 4 && !(player1.Score+player2.Score == 6) {
		s := baseScore[player1.Score]

		if player1.Score == player2.Score {
			result = s + "-All"
		} else {
			result = s + "-" + baseScore[player2.Score]
		}
	}

	if player1.Score == player2.Score {
		result = "Deuce: " + result
	}

	return result, nil
}

func NewPlayerService(repo interfaces.IPlayerRepository) *playerService {
	return &playerService{
		repo: repo,
	}
}
