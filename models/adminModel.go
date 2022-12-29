package models

type Admin struct {
	ID       uint   `json:"Id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Phone    string `json:"Phone"`
}
