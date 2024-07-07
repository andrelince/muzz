package entity

import "time"

type UserInput struct {
	Email    string
	Password string
	Name     string
	Gender   string
	DOB      string
}

type User struct {
	ID       int64
	Email    string
	Password string
	Name     string
	Gender   string
	DOB      time.Time
}

type Token struct {
	Token   string
	Expires int64
}

type Match struct {
	ID      int
	User1ID int
	User2ID int
	IsMatch bool
}
