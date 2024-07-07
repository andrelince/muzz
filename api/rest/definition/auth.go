package definition

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}
