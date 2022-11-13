package render

import (
	"ginEssential/model"
)

func BuildBook(tag *model.MusicBook) *model.MusicBookVO {
	if tag == nil {
		return nil
	}
	bookvo := model.MusicBookVO{
		BookTitle: tag.BookTitle,
		Author:    tag.Author,
		MusicKey:  tag.MusicKey,
		MusicTime: tag.MusicTime,
		Singer:    tag.Singer,
	}
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
