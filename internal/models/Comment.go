package models

import "time"

type Comment struct {
	ID            int       `json:"id"`
	ParentComment *int      `json:"parent_comment"`
	ThreadID      int       `json:"thread_id"`
	UserID        int       `json:"user_id"`
	Body          string    `json:"body"`
	CreatedAT     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
