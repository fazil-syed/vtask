package schemas

type UserRegisterInputSchema struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserLoginInputSchema struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserProfileDataResponse struct {
	Email    *string `json:"email"`
	UserName *string `json:"user_name"`
}

type IdentityInput struct {
	Issuer        string
	Subject       string
	Email         *string
	EmailVerified bool
	PasswordHash  *string // only for password auth
	//Optional
	UserName *string
}
