package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	Registration time.Time `json:"registration"`
	Status       string    `json:"status"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Categories  []Category `json:"categories"`
	Authors     []User     `json:"authors"`
	Feedbacks   []Feedback `json:"feedbacks"`
	Votes       int64      `json:"votes"`
	CreateDate  time.Time  `json:"create_date"`
}

type Feedback struct {
	ID         int64     `json:"id"`
	User       User      `json:"user"`
	Text       string    `json:"text"`
	CreateDate time.Time `json:"create_date"`
}

type Vote struct {
	ID      int64   `json:"id"`
	User    User    `json:"user"`
	Project Project `json:"project"`
}
