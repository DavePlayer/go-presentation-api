package models

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserToUpdate struct {
	Id   int  `json:"id"`
	User User `json:"User"`
}
