package models

import "time"

type User struct {
	Name       string
	Email      string
	HashedPass string // sihhh, it's secret
}

var dbUser map[string]*User

func GetDBUser() map[string]*User {
	if dbUser == nil {
		dbUser = map[string]*User{}
	}
	return dbUser
}

type Post struct {
	UUID    string
	User    string
	Text    string
	Created time.Time
}

var dbPost map[string]*Post

func GetDBPost() map[string]*Post {
	if dbPost == nil {
		dbPost = map[string]*Post{}
	}
	return dbPost
}
