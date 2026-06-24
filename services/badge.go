package services

import (
	"chatter-wails/internal/api"
	"context"
	"log"
	"slices"
)

type BadgeService struct {
	Ctx context.Context

	GlobalBadgeSets *[]api.ApiBadgeSet
}

func NewBadgeService() *BadgeService {
	return &BadgeService{}
}


func (bs *BadgeService) GetGlobalBadgeSets(accessToken string) (*[]api.ApiBadgeSet, error) {
	res, err := api.ApiGetGlobalBadges(accessToken)
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

func GetChannelBadgeSets(accessToken string, broadcasterId string) (*[]api.ApiBadgeSet, error) {
	res, err := api.ApiGetChannelBadges(accessToken, map[string][]string{
		"broadcaster_id": {broadcasterId},
	})
	if err != nil {
		log.Printf("[GetChannelBadgeSets]: An error occurred fetching channel badge sets, aborting\n%+v\n\n", err)
		return nil, err
	}
	if res.Status != 200 {
		log.Printf("[GetChannelBadges]: Failed to get channel badge sets, aborting\n%+v\n\n", res)
		return nil, &api.StatusError[api.ApiGetChannelBadgesRes]{ Res: res}
	}

	return &res.Body.Data, nil
}

func combineBadgeSets(set1 *api.ApiBadgeSet, set2 *api.ApiBadgeSet) *api.ApiBadgeSet {
    if set1 == nil || set2 == nil || set1.Set_id != set2.Set_id {
		return nil
	}

	newVersions := slices.Concat(set1.Versions, set2.Versions)
    return &api.ApiBadgeSet{
        Set_id: set1.Set_id,
        Versions: newVersions,
    }
}


func CombineChannelGlobalSets (cBadgeSets *[]api.ApiBadgeSet,  gBadgeSets *[]api.ApiBadgeSet) *[]api.ApiBadgeSet {
	if cBadgeSets == nil && gBadgeSets == nil {
		return &[]api.ApiBadgeSet{}
	} else if gBadgeSets == nil {
		return cBadgeSets
	} else if cBadgeSets == nil {
		return gBadgeSets
	}


	cSubSetIndex := slices.IndexFunc(*cBadgeSets, func(s api.ApiBadgeSet) bool {
		return s.Set_id == "subscriber"
	})
	gSubSetIndex := slices.IndexFunc(*gBadgeSets, func(s api.ApiBadgeSet) bool {
		return s.Set_id == "subscriber"
	})
	
	var cSubSet *api.ApiBadgeSet = nil
	var gSubSet *api.ApiBadgeSet = nil

	if cSubSetIndex != -1 {
		bs := *cBadgeSets
		cSubSet = &bs[cSubSetIndex]
	}
	if gSubSetIndex != -1 {
		bs := *gBadgeSets
		gSubSet = &bs[gSubSetIndex]
	}
	

	subscriberSet := combineBadgeSets(cSubSet, gSubSet);

	combined := slices.Concat(*cBadgeSets, *gBadgeSets)

    if(subscriberSet != nil) {
		subIndex := slices.IndexFunc(combined, func(s api.ApiBadgeSet) bool {
			return s.Set_id == "subscriber"
		})

        combined[subIndex] = *subscriberSet;
    }

    return &combined;
}
