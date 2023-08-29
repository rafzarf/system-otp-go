package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	Telegram_chat_id string `gorm:"type:varchar(32)"`

	Otp_enabled  bool `gorm:"default:false;"`
	Otp_verified bool `gorm:"default:false;"`

	Otp_secret   string
	Otp_auth_url string
}

func (user *User) BeforeCreate(*gorm.DB) error {
	newUUID := uuid.NewV4()
	user.ID = newUUID.String()

	return nil
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type SendOTPTelegramInput struct {
	UserID         string `json:"user_id" binding:"required"`
	TelegramChatID string `json:"telegram_chat_id" binding:"required"`
}
