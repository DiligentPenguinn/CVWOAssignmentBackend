package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllThreads() ([]*models.Thread, error)
	SingleThread(id int) (*models.Thread, error)
}
