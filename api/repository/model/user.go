package model

import "time"

type UserInput struct {
	Email    string `db:"email"`
	Password string `db:"password"`
	Name     string `db:"name"`
	Gender   string `db:"gender"`
	DOB      string `db:"date_of_birth"`
}

type User struct {
	ID       int64     `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Name     string    `db:"name"`
	Gender   string    `db:"gender"`
	DOB      time.Time `db:"date_of_birth"`
}
