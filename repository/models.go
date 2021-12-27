package repository

import (
	"time"
)

type User struct {
	Id           int64
	TgId         int64
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	IsBot        bool
	Created      time.Time
}

type Message struct {
	Datetime time.Time
	Text     string
	UserID   uint
	ChatID   uint
}

type Chat struct {
	ChatId int64
	Title  string
}
