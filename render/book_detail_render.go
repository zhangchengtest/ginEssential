package render

import (
	"ginEssential/model"
	"strconv"
)

func BuildBookDetail(tag *model.BookDetail) *model.BookDetailVO {
	if tag == nil {
		return nil
	}
	bookvo := model.BookDetailVO{
		Id:          strconv.FormatInt(tag.Id, 10),
		BookContent: tag.BookContent,
		Lyric:       tag.Lyric,
	}
	return &bookvo
}

func BuildBookDetails(tags []model.BookDetail) []model.BookDetailVO {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.BookDetailVO
	for _, tag := range tags {
		responses = append(responses, *BuildBookDetail(&tag))
	}
	return responses
}
