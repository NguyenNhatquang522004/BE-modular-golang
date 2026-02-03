package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/delivery/res"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/IRepositoryPostgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/infrastructure/repository/keycloak"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/utils"
	"github.com/google/uuid"
)

type IUserAuthService interface {
	Login(email string, password string) (*response.Response, error)
	RegisterOne(email string) (*response.Response, error)
	RegisterTwo(email string, otp string) (*response.Response, error)
	RegisterThree(email string, username string, password string) (*response.Response, error)
	ReSendOTP(email string) (*response.Response, error)
	SendLinkResetPassword(email string) (*response.Response, error)
	LogOut(userID string) (*response.Response, error)
	ResetPassword(email string, newPassword string) (*response.Response, error)
}

type UserAuthUseCase struct {
	userRepo       IRepositoryPostgres.IUserRepository
	keycloakClient *keycloak.KeycloakRepository
}

func NewUserAuthUseCase(userRepo IRepositoryPostgres.IUserRepository, keycloakClient *keycloak.KeycloakRepository) *UserAuthUseCase {
	return &UserAuthUseCase{
		userRepo:       userRepo,
		keycloakClient: keycloakClient,
	}
}

func (u *UserAuthUseCase) Login(email string, password string, code string, redirectURI string) (*response.Response, error) {
	if email == "" || password == "" {
		return response.NewResponse(
			response.WithMessage("Email and password must not be empty"),
			response.WithStatus("400"),
		), errors.New("email and password must not be empty")
	}
	checkemail, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return response.NewResponse(
			response.WithMessage("Invalid email "),
			response.WithStatus("401"),
		), errors.New("invalid email ")
	}

	checkpassword := utils.CheckPasswordHash(password, checkemail.Password)
	if !checkpassword {
		return response.NewResponse(
			response.WithMessage("Invalid password"),
			response.WithStatus("401"),
		), errors.New("invalid password")
	}

	// 2. Gọi Keycloak để lấy Token
	tokenResult, err := u.keycloakClient.LoginWithPassword(email, checkemail.Password)
	if err != nil {
		return response.NewResponse(
			response.WithMessage("Error logging in to Keycloak"),
			response.WithStatus("500"),
		), errors.New("error logging in to Keycloak")
	}

	// Giải mã token để lấy thông tin user
	claims, err := u.keycloakClient.DecodeAccessToken(tokenResult.AccessToken)
	if err != nil {
		return nil, err
	}
	mapClaims := *claims
	sub, _ := mapClaims["sub"].(string)

	// 3. Kiểm tra User trong DB
	if checkemail.KeycloakID != sub {
		return response.NewResponse(response.WithData(""),
			response.WithMessage("user keycloak ID does not match."),
			response.WithStatus("404")), errors.New("user keycloak ID does not match.")
	}
	return response.NewResponse(
		response.WithData(&res.TokenResponse{
			AccessToken:  tokenResult.AccessToken,
			RefreshToken: tokenResult.RefreshToken,
			ExpiresIn:    tokenResult.ExpiresIn,
		}),
		response.WithMessage("User logged in successfully"),
		response.WithStatus("200"),
	), nil
}
func (u *UserAuthUseCase) RegisterOne(email string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), errors.New("error checking email")
	}
	if checkEmail != nil {
		return response.NewResponse(
			response.WithMessage("Email already exists"),
			response.WithStatus("400"),
		), errors.New("email already exists")
	}
	user := &entity.User{
		Email:        email,
		StepRegister: 1,
		OTPCode:      uuid.NewString()[:10],
		OTPExpiry:    time.Now().Add(30 * time.Minute),
	}
	_, createUserErr := u.userRepo.CreateUser(user)
	if createUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error creating user"),
			response.WithStatus("500"),
		), errors.New("error creating user")
	}
	return response.NewResponse(
		response.WithMessage("Email is available for registration"),
		response.WithStatus("200"),
		response.WithData("Available"),
	), nil
}

func (u *UserAuthUseCase) RegisterTwo(email string, otp string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), nil
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), nil
	}
	if checkEmail.OTPExpiry.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("OTP has expired"),
			response.WithStatus("400"),
		), nil
	}
	if checkEmail.OTPCode != otp {
		checkEmail.OTPAttempts += 1
		_ = u.userRepo.UpdateUser(checkEmail)
		return response.NewResponse(
			response.WithMessage("Invalid OTP code"),
			response.WithStatus("401"),
		), nil
	}

	checkEmail.StepRegister = 2
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user"),
			response.WithStatus("500"),
		), nil
	}
	return response.NewResponse(
		response.WithMessage("OTP verified successfully"),
		response.WithStatus("200"),
		response.WithData("OTP Verified"),
	), nil
}

func (u *UserAuthUseCase) RegisterThree(email string, username string, password string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), errors.New("error checking email")
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), errors.New("email does not exist")
	}
	if checkEmail.StepRegister != 2 {
		return response.NewResponse(
			response.WithMessage("Previous registration steps not completed"),
			response.WithStatus("400"),
		), errors.New("previous registration steps not completed")
	}
	passwordhash, checkPasswordErr := utils.HashPassword(password)
	if checkPasswordErr != nil {
		return response.NewResponse(
			response.WithMessage("Error hashing password"),
			response.WithStatus("500"),
		), errors.New("error hashing password")

	}

	checkEmail.Username = username
	checkEmail.Password = passwordhash
	checkEmail.IsActive = true
	checkEmail.StepRegister = 3
	enabled := true
	userKC := gocloak.User{
		Username:      &checkEmail.Username,
		Email:         &checkEmail.Email,
		Enabled:       &enabled,
		EmailVerified: &enabled,
	}
	keycloakID, err := u.keycloakClient.CreateUser(&userKC, passwordhash)
	if err != nil {
		return response.NewResponse(
			response.WithMessage("Error creating user in Keycloak"),
			response.WithStatus("500"),
		), errors.New("error creating user in Keycloak")
	}
	checkEmail.KeycloakID = keycloakID
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user"),
			response.WithStatus("500"),
		), errors.New("error updating user")
	}
	return response.NewResponse(
		response.WithMessage("Registration step three not implemented yet"),
		response.WithStatus("501"),
	), nil
}
func (u *UserAuthUseCase) CheckResendOTP(user *entity.User) (*response.Response, error) {

	if !user.OTPTimeWaitOTP.IsZero() && user.OTPTimeWaitOTP.After(time.Now()) {
		return response.NewResponse(
			response.WithMessage("Please wait before requesting a new OTP"),
			response.WithStatus("429"),
		), errors.New("please wait before requesting a new OTP")
	}
	return response.NewResponse(
		response.WithMessage("OTP can be resent"),
		response.WithStatus("200"),
	), nil
}
func (u *UserAuthUseCase) ReSendOTP(email string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), errors.New("error checking email")
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), errors.New("email does not exist")
	}
	if checkEmail.StepRegister != 1 {
		return response.NewResponse(
			response.WithMessage("User is not in step 1 of registration"),
			response.WithStatus("400"),
		), errors.New("user is not in step 1 of registration")
	}
	checkresend, waitOTPRespErr := u.CheckResendOTP(checkEmail)
	if waitOTPRespErr != nil {
		return nil, waitOTPRespErr
	}
	if checkresend.Status != "200" {
		return checkresend, nil
	}
	// Generate new OTP
	checkEmail.OTPCode = uuid.NewString()[:10]
	checkEmail.OTPExpiry = time.Now().Add(30 * time.Minute)
	checkEmail.OTPTimeWaitOTP = time.Now().Add(30 * time.Minute)
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user with new OTP"),
			response.WithStatus("500"),
		), errors.New("error updating user with new OTP")
	}
	// call API to send OTP here notification  module
	return response.NewResponse(
		response.WithMessage("OTP sent successfully"),
		response.WithStatus("200"),
	), nil
}
func (u *UserAuthUseCase) SendLinkResetPassword(email string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), errors.New("error checking email")
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), errors.New("email does not exist")
	}
	checkresend, waitOTPRespErr := u.CheckResendOTP(checkEmail)
	if waitOTPRespErr != nil {
		return nil, waitOTPRespErr
	}
	if checkresend.Status != "200" {
		return checkresend, nil
	}
	// Generate new OTP for reset password
	checkEmail.OTPCode = uuid.NewString()[:10]
	checkEmail.OTPExpiry = time.Now().Add(30 * time.Minute)
	checkEmail.OTPTimeWaitOTP = time.Now().Add(30 * time.Minute)
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user with new OTP"),
			response.WithStatus("500"),
		), errors.New("error updating user with new OTP")
	}
	// call API to send OTP here notification  module
	return response.NewResponse(
		response.WithMessage("OTP sent successfully"),
		response.WithStatus("200"),
	), nil
}
func (u *UserAuthUseCase) ResetPassword(email string, newPassword string) (*response.Response, error) {
	checkEmail, checkEmailErr := u.userRepo.GetUserByEmail(email)
	if checkEmailErr != nil {
		return response.NewResponse(
			response.WithMessage("Error checking email"),
			response.WithStatus("500"),
		), errors.New("error checking email")
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), errors.New("email does not exist")
	}
	hashedPassword, hashErr := utils.HashPassword(newPassword)
	if hashErr != nil {
		return response.NewResponse(
			response.WithMessage("Error hashing new password"),
			response.WithStatus("500"),
		), errors.New("error hashing new password")
	}
	checkEmail.Password = hashedPassword
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user password"),
			response.WithStatus("500"),
		), errors.New("error updating user password")
	}
	return response.NewResponse(
		response.WithMessage("Password reset successfully"),
		response.WithStatus("200"),
	), nil
}

func (u *UserAuthUseCase) LogOut(userID string, accessToken string) (*response.Response, error) {
	user, err := u.userRepo.GetUserByID(uuid.MustParse(userID))
	if err != nil {
		return response.NewResponse(
			response.WithMessage("Error retrieving user"),
			response.WithStatus("500"),
		), errors.New("error retrieving user")
	}
	if user == nil {
		return response.NewResponse(
			response.WithMessage("User not found"),
			response.WithStatus("404"),
		), errors.New("user not found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = u.keycloakClient.Logout(ctx, accessToken)
	if err != nil {
		return response.NewResponse(
			response.WithMessage("Error logging out user from Keycloak"),
			response.WithStatus("500"),
		), errors.New("error logging out user from Keycloak")
	}
	return response.NewResponse(
		response.WithMessage("User logged out successfully"),
		response.WithStatus("200"),
	), nil
}
