package model

type MusicBookDTO struct {
	PageInfo
	BookTitle string
	Author    string
}

type BookDetailDTO struct {
	PageInfo
	Id          string
	BookId      string
	BookContent string
	Lyric       string
}

type BookPieceDTO struct {
	BookId      string
	BookContent string
	Connection  string
	UpPoints    string
	Line1       string
	Line2       string
	DownPoints  string
	Lyric       string
	Indent      string
	BookOrder   int `json:"bookOrder"`
	PhaseId     string
	BreakFlag   int
}

type PieceContentDTO struct {
	BookId    string
	PhaseId   string
	BookOrder int `json:"bookOrder"`
}

type FileTemp2FormalDTO struct {
	ServeCode    string
	RelativePath string
}

type PuzzlePieceDTO struct {
	Url string `json:"url"`
}
