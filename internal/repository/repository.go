package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllThreads() ([]*models.Thread, error)
	SingleThread(id int) (*models.Thread, error)
	GetCommentsByThreadID(id int) ([]*models.Comment, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	InsertThread(thread models.Thread) (int, error)
	InsertComment(comment models.Comment) (int, error)
	InsertReply(reply models.Reply) (int, error)
}
