package controllers

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pquerna/otp/totp"
	"github.com/rafzarf/system-otp-go/helper"
	"github.com/rafzarf/system-otp-go/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Hash the input password using MD5
	hashedPassword := helper.HashMD5(payload.Password)

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: hashedPassword, // Store the hashed password
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Email already exist, please use another email address"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Registered successfully, please login"})
}

func (ac *AuthController) LoginUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	// Hash the input password using MD5 for comparison
	hashedInputPassword := helper.HashMD5(payload.Password)

	// Use the helper function to compare passwords
	if !helper.CompareHashedPasswords(ctx, hashedInputPassword, user) {
		return
	}

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": userResponse})
}

func (ac *AuthController) GenerateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Retrieve user from the database
	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User doesn't exist"})
		return
	}

	// Generate TOTP token
	token, err := totp.GenerateCode(user.Otp_secret, time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate TOTP token"})
		return
	}

	// Send the TOTP token in the response
	ctx.JSON(http.StatusOK, gin.H{"otp_token": token})
}

func (ac *AuthController) SendOTPTelegram(ctx *gin.Context) {
	var payload *models.SendOTPTelegramInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User doesn't exist"})
		return
	}

	// Generate OTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "contoh-sistem-otp",
		AccountName: user.Email, // Use the user's email or some identifier
		SecretSize:  6,
		Period:      180,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate OTP"})
		return
	}

	// Update user's OTP secret and URL
	dataToUpdate := models.User{
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	ac.DB.Model(&user).Updates(dataToUpdate)

	//bot mulai
	bot, err := tgbotapi.NewBotAPI("6579966415:AAGs4-qcypUQOS51rf5d5WNxT3Wf1mDpMlQ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to initialize Telegram bot"})
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here is your OTP information:\n"+
				"OTP Secret: "+key.Secret()+"\n"+
				"OTP Auth URL: "+key.URL())
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (ac *AuthController) VerifyOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Token is invalid or user doesn't exist"

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	dataToUpdate := models.User{
		Otp_enabled:  true,
		Otp_verified: true,
	}

	ac.DB.Model(&user).Updates(dataToUpdate)

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"otp_verified": true, "user": userResponse})
}

func (ac *AuthController) ValidateOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "Token is invalid or user doesn't exist"

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"otp_valid": true})
}

func (ac *AuthController) DisableOTP(ctx *gin.Context) {
	var payload *models.OTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User doesn't exist"})
		return
	}

	user.Otp_enabled = false
	ac.DB.Save(&user)

	userResponse := gin.H{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	ctx.JSON(http.StatusOK, gin.H{"otp_disabled": true, "user": userResponse})
}
