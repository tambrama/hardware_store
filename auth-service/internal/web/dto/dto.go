package dto

type LoginInput struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
    AppID    string `validate:"required,gt=0"`
}

type LoginOutput struct {
    AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterInput struct {
	Email       string `validate:"required,email"`
	Password    string `validate:"required,min=8,max=64"` 
	Name        string `validate:"required,min=2,max=50"`
	Surname     string `validate:"required,min=2,max=50"`
	PhoneNumber string `validate:"required,e164"` 
}

type ChangePasswordInput struct {
	Email       string `validate:"required,email"`
	OldPassword string `validate:"required,min=8"`
	NewPassword string `validate:"required,min=8"`
}