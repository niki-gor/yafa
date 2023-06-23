package model

import "time"

type Forum struct {
	Id      int64 `json:"-"`
	Title   string
	User    string
	Slug    string
	Posts   int
	Threads int
}
type Post struct {
	Id       int64
	Parent   int64
	Author   string
	Message  string
	IsEdited bool `json:"isEdited"`
	Forum    string
	Thread   int32
	Created  time.Time
	Path     int64 `json:"-"`
}
type Posts struct {
	Posts []Post
}
type PostAll struct {
	Post   *Post
	Author *User
	Thread *Thread
	Forum  *Forum
}
type Status struct {
	User   int
	Forum  int
	Thread int
	Post   int
}
type Thread struct {
	Id      int
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int
	Slug    string
	Created time.Time
}
type User struct {
	Id       int `json:"-"`
	Nickname string
	Fullname string
	About    string
	Email    string
}
type Users struct {
	Users []User
}
type VoteRequest struct {
	Nickname string
	Voice    int
}
type Vote struct {
	Id     int `json:"-"`
	User   int
	Thread int
	Voice  int
}
