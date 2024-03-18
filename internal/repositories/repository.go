package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Armunz/learn-docker/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, user entity.User) error
	Get(ctx context.Context, limit int, offset int) ([]entity.User, error)
	GetById(ctx context.Context, userID string) (entity.User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, userID string) error
}

type repositoryImpl struct {
	db        *sql.DB
	timeoutMs int
}

func NewRepository(db *sql.DB, timeoutMs int) Repository {
	return &repositoryImpl{
		db:        db,
		timeoutMs: timeoutMs,
	}
}

// Count implements Repository.
func (r *repositoryImpl) Count(ctx context.Context) (int64, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	query := `SELECT COUNT(*) FROM users`

	var result int64
	err := r.db.QueryRowContext(ctxTimeout, query).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Create implements Repository.
func (r *repositoryImpl) Create(ctx context.Context, user entity.User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	query := `INSERT INTO users(user_id, name, age) VALUES(?,?,?)`
	_, err := r.db.ExecContext(ctxTimeout, query, user.UserID, user.Name, user.Age)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements Repository.
func (r *repositoryImpl) Delete(ctx context.Context, userID string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	query := `DELETE FROM users WHERE user_id = ?`
	_, err := r.db.ExecContext(ctxTimeout, query, userID)
	if err != nil {
		return err
	}

	return nil
}

// Get implements Repository.
func (r *repositoryImpl) Get(ctx context.Context, limit int, offset int) ([]entity.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	query := `SELECT user_id, name, age FROM users LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctxTimeout, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Age); err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetById implements Repository.
func (r *repositoryImpl) GetById(ctx context.Context, userID string) (entity.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	var user entity.User
	query := `SELECT user_id, name, age FROM users WHERE user_id = ?`
	if err := r.db.QueryRowContext(ctxTimeout, query, userID).Scan(&user.UserID, &user.Name, &user.Age); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// Update implements Repository.
func (r *repositoryImpl) Update(ctx context.Context, user entity.User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	query := `UPDATE users SET name = ?, age = ? WHERE user_id = ?`
	_, err := r.db.ExecContext(ctxTimeout, query, user.Name, user.Age, user.UserID)

	return err
}
