package todo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, title string) (Todo, error) {
	id := uuid.NewString()
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (id, title, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		id, title, false, now, now,
	)
	if err != nil {
		return Todo{}, err
	}
	return Todo{ID: id, Title: title, Completed: false, CreatedAt: now, UpdatedAt: now}, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Todo, error) {
	var t Todo
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, completed, created_at, updated_at FROM todos WHERE id=$1`, id,
	)
	if err := row.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, err
	}
	return t, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, completed, created_at, updated_at FROM todos ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, update UpdateTodoRequest) (Todo, error) {
	// Fetch current
	current, err := r.Get(ctx, id)
	if err != nil {
		return Todo{}, err
	}
	if update.Title != nil {
		current.Title = *update.Title
	}
	if update.Completed != nil {
		current.Completed = *update.Completed
	}
	current.UpdatedAt = time.Now().UTC()
	_, err = r.db.ExecContext(ctx,
		`UPDATE todos SET title=$1, completed=$2, updated_at=$3 WHERE id=$4`,
		current.Title, current.Completed, current.UpdatedAt, id,
	)
	if err != nil {
		return Todo{}, err
	}
	return current, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
