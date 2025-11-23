package dto

type UserCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
	Number   string `json:"number"`
	Role     string `json:"role"`
}
