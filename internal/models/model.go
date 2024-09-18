package models

type User struct {
	Id       int
	Username string
	Projects []Project
}

type Project struct {
	Id         int
	Name       string
	Desciption string
	LikesCount int
	Comments   []Comment
}

type Comment struct {
	Id     int
	UserId int
}
