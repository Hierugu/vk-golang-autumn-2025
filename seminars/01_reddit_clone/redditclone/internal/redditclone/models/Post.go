package models

type Post struct {
	Id       string
	Title    string
	Content  string
	Author   *User
	Comments []Comment
}

type Comment struct {
	Id      string
	Content string
	Author  *User
}
