package render

import (
	"ginEssential/model"
	"ginEssential/util"
)

func BuildBook(tag *model.MusicBook) *model.MusicBookVO {
	if tag == nil {
		return nil
	}
	bookvo := model.MusicBookVO{}
	util.SimpleCopyProperties(&bookvo, &tag)
	return &bookvo
}

func BuildBooks(tags []model.MusicBook) *[]model.MusicBookVO {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.MusicBookVO
	for _, tag := range tags {
		responses = append(responses, *BuildBook(&tag))
	}
	return &responses
}
