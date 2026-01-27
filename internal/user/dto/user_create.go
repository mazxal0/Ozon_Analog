package dto

type UserChange struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	LastName     string `json:"last_name"`
	Number       string `json:"number"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
}

type PasswordChange struct {
	Password string `json:"password"`
}
