package render

import (
	"ginEssential/model"
	"ginEssential/util"
)

func BuildFood(tag *model.Food) *model.FoodVO {
	if tag == nil {
		return nil
	}
	articleVO := model.FoodVO{}

	util.SimpleCopyProperties(&articleVO, tag)
	return &articleVO
}

func BuildFoods(tags []model.Food) *[]model.FoodVO {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.FoodVO
	for _, tag := range tags {
		responses = append(responses, *BuildFood(&tag))
	}
	return &responses
}
