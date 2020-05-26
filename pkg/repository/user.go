package repository

import (
	"context"
	"database/sql"

	"github.com/PedPet/user/model"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

const (
	// InsertUser is a sql statement to insert a user into the users database
	InsertUser string = "INSERT INTO users (username) VALUES(?)"
	// GetUser is a sql statement to get a user from the users database
	GetUser string = "SELECT id FROM users WHERE username = ?"
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
	logger := log.With(r.logger, "method", "CreateUser")
	// logg.Println("CreateUser", ctx, user, r.db)
	stmt, err := r.db.PrepareContext(ctx, InsertUser)
	if err != nil {
		return errors.Wrap(err, "Failed to prepare insert user statement")
	}
	defer stmt.Close()

	if user.Username == "" {
		return errRepo
	}

	result, err := stmt.Exec(user.Username)
	if err != nil {
		return errors.Wrap(err, "Failed to execute prepared insert statement")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Failed to get last insert id")
	}

	user.ID = int(id)
	logger.Log("Create user", user.ID)

	return nil
}

func (r repo) GetUser(ctx context.Context, user *model.User) error {
	row, err := r.db.Query(GetUser, user.Username)
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

	row.Scan(&user.ID)
	return nil
}
