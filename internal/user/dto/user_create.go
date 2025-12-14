package dto

type UserChange struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	LastName string `json:"last_name"`
	Number   string `json:"number"`
}

type PasswordChange struct {
	Password string `json:"password"`
}
