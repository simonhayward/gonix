package gonix

import (
	"time"
)

type User struct {
	UserID  int       `json:"user_id"`
	Created time.Time `json:"created"`
	Name    string    `json:"name"`
}

type Token struct {
	Token    string    `json:"token"`
	UserID   int       `json:"user_id"`
	Created  time.Time `json:"created"`
	LastUsed time.Time `json:"last_used"`
}
