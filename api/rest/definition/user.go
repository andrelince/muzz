package definition

type UserInput struct {
	Email        string   `json:"email" validate:"required,email"`
	Password     string   `json:"password" validate:"required"` // implement hash
	Name         string   `json:"name" validate:"required"`
	Gender       string   `json:"gender" validate:"oneof=M F"`
	DOB          string   `json:"dob" validate:"required,dob"`
	LocationLat  *float64 `json:"locationLat"`
	LocationLong *float64 `json:"locationLong"`
}

type User struct {
	ID           int64    `json:"id"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	Name         string   `json:"name"`
	Gender       string   `json:"gender"`
	Age          int      `json:"age"`
	LocationLat  *float64 `json:"location_lat,omitempty"`
	LocationLong *float64 `json:"location_long,omitempty"`
}

type SwipeInput struct {
	UserID     int    `json:"user_id" validate:"required"`
	Preference string `json:"preference" validate:"oneof=yes no"`
}

type Match struct {
	MatchID *int `json:"match_id,omitempty"`
	Matched bool `json:"matched"`
}

type Discovery struct {
	User                User    `json:"user"`
	DistanceFromMe      float64 `json:"distance"`
	AttractivenessScore int     `json:"attractiveness"`
}
