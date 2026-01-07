package services

import (
	"chatter-wails/internal/api"
	"context"
	"log"
)

type BadgeService struct {
	Ctx context.Context

	GlobalBadgeSets *[]api.ApiBadgeSet
}

func NewBadgeService() *BadgeService {
	return &BadgeService{}
}


func (bs *BadgeService) GetGlobalBadgeSets(access_token string) (*[]api.ApiBadgeSet, error) {
	res, err := api.ApiGetGlobalBadges(access_token)
	if err != nil {
		log.Printf("[GetGlobalBadgeSets]: An error occurred fetching global badge sets, aborting\n\n")
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetGlobalBadges]: Failed to get global badge sets, aborting\n\n")
		return nil, &api.StatusError[api.ApiGetGlobalBadgesRes]{ Res: res }
	}

	sets := &res.Body.Data

	bs.GlobalBadgeSets = sets

	return sets, nil
}
