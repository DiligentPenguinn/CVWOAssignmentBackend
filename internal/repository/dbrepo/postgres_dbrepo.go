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
			id, user_id, title, body,
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

func (m *PostgresDBRepo) AllThreadsWithUsers() ([]*models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, user_id, title, body,
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

func (m *PostgresDBRepo) SingleThread(id int) (*models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	select 
		id, user_id, title, body, updated_at
	from 
		threads 
	where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var thread models.Thread

	err := row.Scan(
		&thread.ID,
		&thread.User_ID,
		&thread.Title,
		&thread.Body,
		&thread.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &thread, err
}

func (m *PostgresDBRepo) GetCommentsByThreadID(id int) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, user_id, body, parent_id,
			created_at, updated_at
		from
			comments
		where
			parent_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Body,
			&comment.ParentID,
			&comment.CreatedAT,
			&comment.UpdatedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}
	return comments, nil
}
