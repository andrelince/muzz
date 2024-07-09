package model

import "time"

type UserInput struct {
	Email        string   `db:"email"`
	Password     string   `db:"password"`
	Name         string   `db:"name"`
	Gender       string   `db:"gender"`
	DOB          string   `db:"date_of_birth"`
	LocationLat  *float64 `db:"location_lat"`
	LocationLong *float64 `db:"location_long"`
}

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	Name         string    `db:"name"`
	Gender       string    `db:"gender"`
	DOB          time.Time `db:"date_of_birth"`
	LocationLat  *float64  `db:"location_lat"`
	LocationLong *float64  `db:"location_long"`
}

type Swipe struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	SwipedUserID int       `db:"swiped_user_id"`
	SwipeStatus  bool      `db:"swipe_status"`
	CreatedAt    time.Time `db:"created_at"`
}

type Match struct {
	ID        int `db:"id"`
	User1ID   int `db:"user1_id"`
	User2ID   int `db:"user2_id"`
	IsMatch   bool
	CreatedAt time.Time `db:"created_at"`
}

type Discovery struct {
	User                User    `db:"user"`
	DistanceFromMe      float64 `db:"distance_from_me"`
	AttractivenessScore int     `db:"attractiveness_score"`
}
