package definition

type UserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // implement hash
	Name     string `json:"name" validate:"required"`
	Gender   string `json:"gender" validate:"oneof=M F"`
	DOB      string `json:"dob" validate:"required,dob"`
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}
