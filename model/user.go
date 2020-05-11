package model

type User struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Confirmed   bool   `json:"confirmed"`
}
