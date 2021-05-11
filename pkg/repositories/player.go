package repositories

import (
	"github.com/afex/hystrix-go/hystrix"
	"gorm.io/gorm"

	"github.com/likakuli/generic-project-template/pkg/interfaces"
	"github.com/likakuli/generic-project-template/pkg/models"
)

type PlayerRepository struct {
	db *gorm.DB
}

func (repository *PlayerRepository) GetPlayerByName(name string) (*models.Player, error) {
	var player models.Player

	err := repository.db.Where("name = ?", name).Find(&player).Error
	if err != nil {
		return nil, err
	}

	return &player, nil
}

type PlayerRepositoryWithCircuitBreaker struct {
	PlayerRepository interfaces.IPlayerRepository
}

func (repository *PlayerRepositoryWithCircuitBreaker) GetPlayerByName(name string) (*models.Player, error) {
	output := make(chan *models.Player, 1)
	hystrix.ConfigureCommand("get_player_by_name", hystrix.CommandConfig{Timeout: 1000})
	errors := hystrix.Go("get_player_by_name", func() error {
		player, _ := repository.PlayerRepository.GetPlayerByName(name)

		output <- player
		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}

func GetDefaultPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}
