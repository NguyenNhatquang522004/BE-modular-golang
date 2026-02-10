package req

type RegisterOneRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterTwoRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type RegisterThreeRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type ReSendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SendLinkResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=50"`
}

type Logout struct {
	UserID      string `json:"user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

type GoogleLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type GoogleLoginURLRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required,url"`
}
