package models

import (
	"time"
)

type User struct {
	User_id        string    `gorm:"type:char(36);primary_key" json:"user_id"`
	User_username  string    `gorm:"unique" json:"user_username"`
	User_name      string    `json:"user_name"`
	User_email     string    `gorm:"unique" json:"user_email"`
	User_gender    string    `json:"user_gender"`
	User_password  string    `json:"user_password"`
	User_role      string    `json:"user_role"`
	User_telp      string    `json:"user_telp"`
	User_birthdate string    `json:"user_birthdate"`
	User_profpic   string    `json:"user_profpic"`
	User_otp       string    `json:"user_otp"`
	User_otp_valid time.Time `json:"user_otp_valid"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
