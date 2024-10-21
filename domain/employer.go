package domain

type Employer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Major   string `json:"major"`
	Image   string `json:"image"`
}
