package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ParentID  int       `json:"parent_id"`
	UserID    int       `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
