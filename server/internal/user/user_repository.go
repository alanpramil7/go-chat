package user

import (
	"context"
	"database/sql"
	"errors"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DB
}

func (r *repository) UserExists(ctx context.Context, username, email string) bool {
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	_ = r.db.QueryRowContext(ctx, checkQuery, username, email).Scan(&exists)

	return exists
}

func NewRepository(db DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	exists := r.UserExists(ctx, user.Username, user.Email)
	if exists {
		return nil, errors.New("No account found")
	}

	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}

	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return &User{}, err
	}

	return &u, nil
}
