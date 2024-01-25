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
	GetRepliesByCommentID(id int) ([]*models.Reply, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	InsertThread(thread models.Thread) (int, error)
	InsertComment(comment models.Comment) error
	InsertReply(reply models.Reply) error
	InsertUser(user models.User) error
}
