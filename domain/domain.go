package domain

type Banner struct {
	ID        uint   `json:"id"`
	Image     string `json:"image"`
	Link      string `json:"link"`
	Is_Active *bool  `json:"is_active"`
}
