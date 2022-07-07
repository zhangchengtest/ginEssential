package vo

	type ArticleVO struct {
	Id int64 `json:"id"`
	Chapter      int32`json:"chapter"`
	Title      string `json:"title"`
	Content      string `json:"content"`
	ReadCount      int32 `json:"readCount"`
	Question string `json:"question"`

}

