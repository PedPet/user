package repository

import (
	"context"
	"database/sql"

	"github.com/PedPet/user/model"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	insertUser string = "INSERT INTO users (username) VALUES($1)"
	getUser    string = "SELECT * FROM users WHERE username = $1"
)

var errRepo = errors.New("Unable to handle Repo Request")

// User interface to define user repo
type User interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, user *model.User) error
}

type repo struct {
	db     *sql.DB
	logger log.Logger
}

// NewRepo creates a new repo instance
func NewRepo(db *sql.DB, logger log.Logger) User {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (r repo) CreateUser(ctx context.Context, user *model.User) error {
	stmt, err := r.db.PrepareContext(ctx, insertUser)
	defer stmt.Close()

	if user.Username == "" {
		return errRepo
	}

	_, err = stmt.ExecContext(ctx, user.Username)
	if err != nil {
		return err
	}

	return nil
}

func (r repo) GetUser(ctx context.Context, user *model.User) error {
	row, err := r.db.Query(getUser, user.Username)
	if err != nil {
		return errors.Wrap(err, "Failed to get user from database")
	}

	err = rowToUser(row, user)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func rowToUser(row *sql.Rows, user *model.User) error {
	if !row.Next() {
		return errors.New("No user found")
	}

	row.Scan(user.ID)
	return nil
}
