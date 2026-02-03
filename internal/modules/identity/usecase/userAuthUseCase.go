package usecase

import (
	"time"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/http/response"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/entity"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/domain/repository_postgres"
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/identity/utils"
	"github.com/google/uuid"
)

type IUserAuthService interface {
	Login(email string, password string) (*response.Response, error)
	RegisterOne(email string) (*response.Response, error)
	RegisterTwo(email string, otp string) (*response.Response, error)
	RegisterThree(email string, username string, password string) (*response.Response, error)
	ReSendOTP(email string) (*response.Response, error)
	ResetPassword(email string, newPassword string) (*response.Response, error)
}

type UserAuthUseCase struct {
	userRepo repository_postgres.IUserRepository
}

func NewUserAuthUseCase(userRepo repository_postgres.IUserRepository) *UserAuthUseCase {
	return &UserAuthUseCase{
		userRepo: userRepo,
	}
}

func (u *UserAuthUseCase) Login(email string, password string) (*response.Response, error) {
	if email == "" || password == "" {
		return response.NewResponse(
			response.WithMessage("Email and password must not be empty"),
			response.WithStatus("400"),
		), nil
	}
	if email, err := u.userRepo.GetUserByEmail(email); err != nil || !utils.CheckPasswordHash(password, email.Password) {
		return response.NewResponse(
			response.WithMessage("Invalid email or password"),
			response.WithStatus("401"),
		), nil
	}
	return response.NewResponse(
		response.WithData("Login successful"),
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
		), nil
	}
	if checkEmail != nil {
		return response.NewResponse(
			response.WithMessage("Email already exists"),
			response.WithStatus("400"),
		), nil
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
		), nil
	}
	return response.NewResponse(
		response.WithMessage("Email is available"),
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
		), nil
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), nil
	}
	if checkEmail.StepRegister != 2 {
		return response.NewResponse(
			response.WithMessage("Previous registration steps not completed"),
			response.WithStatus("400"),
		), nil
	}
	passwordhash, checkPasswordErr := utils.HashPassword(password)
	if checkPasswordErr != nil {
		return response.NewResponse(
			response.WithMessage("Error hashing password"),
			response.WithStatus("500"),
		), nil

	}
	checkEmail.Username = username
	checkEmail.Password = passwordhash
	checkEmail.IsActive = true
	checkEmail.StepRegister = 3
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user"),
			response.WithStatus("500"),
		), nil
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
		), nil
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
		), nil
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), nil
	}
	if checkEmail.StepRegister != 1 {
		return response.NewResponse(
			response.WithMessage("User is not in step 1 of registration"),
			response.WithStatus("400"),
		), nil
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
		), nil
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
		), nil
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), nil
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
		), nil
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
		), nil
	}
	if checkEmail == nil {
		return response.NewResponse(
			response.WithMessage("Email does not exist"),
			response.WithStatus("400"),
		), nil
	}
	hashedPassword, hashErr := utils.HashPassword(newPassword)
	if hashErr != nil {
		return response.NewResponse(
			response.WithMessage("Error hashing new password"),
			response.WithStatus("500"),
		), nil
	}
	checkEmail.Password = hashedPassword
	updateUserErr := u.userRepo.UpdateUser(checkEmail)
	if updateUserErr != nil {
		return response.NewResponse(
			response.WithMessage("Error updating user password"),
			response.WithStatus("500"),
		), nil
	}
	return response.NewResponse(
		response.WithMessage("Password reset successfully"),
		response.WithStatus("200"),
	), nil
}
