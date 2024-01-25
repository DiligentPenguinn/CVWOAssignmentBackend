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
			&thread.UserID,
			&thread.Title,
			&thread.Body,
			&thread.CreatedAt,
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
			&thread.UserID,
			&thread.Title,
			&thread.Body,
			&thread.CreatedAt,
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
		&thread.UserID,
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
			&comment.CreatedAt,
			&comment.UpdatedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}
	return comments, nil
}

func (m *PostgresDBRepo) GetRepliesByCommentID(id int) ([]*models.Reply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, user_id, body, parent_id,
			created_at, updated_at
		from
			replies
		where
			parent_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*models.Reply

	for rows.Next() {
		var reply models.Reply
		err := rows.Scan(
			&reply.ID,
			&reply.UserID,
			&reply.Body,
			&reply.ParentID,
			&reply.CreatedAt,
			&reply.UpdatedAt)
		if err != nil {
			return nil, err
		}

		replies = append(replies, &reply)
	}
	return replies, nil
}

func (m *PostgresDBRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, first_name, last_name, username, email, password,
			created_at, updated_at from users where username = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID returns one use, by ID.
func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, first_name, last_name, username, email, password,
			created_at, updated_at from users where id = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) InsertThread(thread models.Thread) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into threads (user_id, title, body, created_at, updated_at)
			values ($1, $2, $3, $4, $5) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		thread.UserID,
		thread.Title,
		thread.Body,
		thread.CreatedAt,
		thread.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) InsertComment(comment models.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into comments (parent_id, user_id, body, created_at, updated_at)
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		comment.ParentID,
		comment.UserID,
		comment.Body,
		comment.CreatedAt,
		comment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertReply(reply models.Reply) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into replies (parent_id, user_id, body, created_at, updated_at)
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		reply.ParentID,
		reply.UserID,
		reply.Body,
		reply.CreatedAt,
		reply.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into users (first_name, last_name, username, 
			email, password, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
