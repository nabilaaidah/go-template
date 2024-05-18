package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserId    string `json:"userId"`
	User      User   `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE;"`
	Token     string `gorm:"unique;size:255"`
	ExpiresAt time.Time
}
