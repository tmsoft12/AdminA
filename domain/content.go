package domain

type Content struct {
	ID       uint   `json:"id"`
	Position string `json:"position"`
	Content  string `json:"content"`
}
