package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllThreads() ([]*models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, user_id, title, body
			created_at, updated_at
		from
			threads`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*models.Thread

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(
			&thread.ID,
			&thread.User_ID,
			&thread.Title,
			&thread.Body,
			&thread.CreatedAT,
			&thread.UpdatedAt)
		if err != nil {
			return nil, err
		}

		threads = append(threads, &thread)
	}
	return threads, nil
}