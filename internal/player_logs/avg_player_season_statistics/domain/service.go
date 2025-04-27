package domain

import (
	"context"
)

type Service interface {
	GetPlayerData(ctx context.Context, seasonID int, nbaPlayerID int) (*AveragePlayerSeason, error)
	Sync(ctx context.Context, ids []int32) error
}

type ServiceImpl struct {
	playerGameStatisticRepository PlayerGameStatisticRepository
}

type PlayerGameStatisticRepository interface {
	GetCacheAveragePlayerSeason(ctx context.Context, seasonID, playerID int) (*AveragePlayerSeason, error)
	GetAveragePlayersSeason(ctx context.Context, ids []int32) ([]AveragePlayerSeason, error)
	SetCacheAveragePlayerSeason(ctx context.Context, seasonID, playerID int, averagePlayerSeason AveragePlayerSeason) error
}

func NewService(playerGameStatisticRepository PlayerGameStatisticRepository) Service {
	return &ServiceImpl{playerGameStatisticRepository: playerGameStatisticRepository}
}

func (c *ServiceImpl) GetPlayerData(ctx context.Context, seasonID int, nbaPlayerID int) (*AveragePlayerSeason, error) {
	return c.playerGameStatisticRepository.GetCacheAveragePlayerSeason(ctx, seasonID, nbaPlayerID)
}

func (c *ServiceImpl) Sync(ctx context.Context, ids []int32) error {

	averagePlayerSeason, err := c.playerGameStatisticRepository.GetAveragePlayersSeason(ctx, ids)
	if err != nil {
		return err
	}
	for _, a := range averagePlayerSeason {

		if err := c.playerGameStatisticRepository.SetCacheAveragePlayerSeason(ctx, a.SeasonID, a.NbaPlayerID, a); err != nil {
			return err
		}
	}
	return nil
}
