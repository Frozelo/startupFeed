package main

import (
	"fmt"
)

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
	Id   int
	User User
}

func main() {
	newUser := User{
		Id:       0,
		Username: "Frozelo",
		Projects: []Project{},
	}
	newComment := Comment{
		Id:   0,
		User: newUser,
	}
	fmt.Printf("the user info is %v", newUser)
	fmt.Printf("the comment info is %v", newComment)
}
