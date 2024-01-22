package models

import "time"

type Thread struct {
	ID        int       `json:"id"`
	User_ID   int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAT time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
