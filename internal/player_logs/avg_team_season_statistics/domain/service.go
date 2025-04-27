package domain

import "context"

type Service interface {
	GetTeamData(ctx context.Context, seasonID int, nbaTeamID int) (*AverageTeamSeason, error)
	Sync(ctx context.Context, ids []int32) error
}

type ServiceImpl struct {
	teamGameStatisticRepository TeamGameStatisticRepository
}

type TeamGameStatisticRepository interface {
	GetCacheAverageTeamSeason(ctx context.Context, seasonID, teamID int) (*AverageTeamSeason, error)
	GetAverageTeamsSeason(ctx context.Context, ids []int32) ([]AverageTeamSeason, error)
	SetCacheAverageTeamSeason(ctx context.Context, seasonID, teamID int, averageTeamSeason AverageTeamSeason) error
}

func NewService(teamGameStatisticRepository TeamGameStatisticRepository) Service {
	return &ServiceImpl{teamGameStatisticRepository: teamGameStatisticRepository}
}

func (c *ServiceImpl) GetTeamData(ctx context.Context, seasonID int, nbaPlayerID int) (*AverageTeamSeason, error) {
	return c.teamGameStatisticRepository.GetCacheAverageTeamSeason(ctx, seasonID, nbaPlayerID)
}

func (c *ServiceImpl) Sync(ctx context.Context, ids []int32) error {

	averagePlayerSeason, err := c.teamGameStatisticRepository.GetAverageTeamsSeason(ctx, ids)
	if err != nil {
		return err
	}
	for _, a := range averagePlayerSeason {

		if err := c.teamGameStatisticRepository.SetCacheAverageTeamSeason(ctx, a.SeasonID, a.NbaTeamID, a); err != nil {
			return err
		}
	}
	return nil
}
