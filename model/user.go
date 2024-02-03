package model

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"Username"`
}
