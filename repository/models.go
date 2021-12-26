package repository

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	TgId         int64
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	IsBot        bool
	Messages     []Message
}

type Message struct {
	gorm.Model
	Datetime time.Time
	Text     string
	UserID   uint
	ChatID   uint
}

type Chat struct {
	gorm.Model
	ChatId   int64
	Title    string
	Messages []Message
}
