package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Project struct {
	Id         int
	Name       string
	Desciption string
	LikesCount int
}

type Comment struct {
	Id     int
	UserId int
}
