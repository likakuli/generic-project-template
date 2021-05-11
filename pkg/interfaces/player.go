package interfaces

import (
	"context"

	"github.com/likakuli/generic-project-template/pkg/models"
)

type IPlayerService interface {
	GetScores(ctx context.Context, player1Name string, player2Name string) (string, error)
}

type IPlayerRepository interface {
	GetPlayerByName(name string) (*models.Player, error)
}
