package msql

import (
	"context"
	"database/sql"

	"github.com/DexScen/WebBook/backend/auth/internal/domain"
	"github.com/DexScen/WebBook/backend/auth/internal/errors"
	_ "github.com/lib/pq"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Register(ctx context.Context, user *domain.User) error {
	tr, err := u.db.Begin()
	if err != nil {
		return err
	}
	statement, err := tr.Prepare("INSERT INTO users (name, email, role, password, login) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Email, user.Role, user.Password, user.Login)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}

func (u *Users) GetPassword(ctx context.Context, login string) (string, error) {
	tr, err := u.db.Begin()
	if err != nil {
		return "", err
	}
	statement, err := tr.Prepare("SELECT password FROM users WHERE login=$1")
	if err != nil {
		tr.Rollback()
		return "", err
	}
	defer statement.Close()

	var password string
	err = statement.QueryRow(login).Scan(&password)
	if err != nil {
		tr.Rollback()
		if err == sql.ErrNoRows {
			return "", errors.ErrUserNotFound
		}
		return "", err
	}

	if err := tr.Commit(); err != nil {
		return "", err
	}

	return password, nil
}

func (u *Users) GetRole(ctx context.Context, login string) (string, error) {
	tr, err := u.db.Begin()
	if err != nil {
		return "", err
	}
	statement, err := tr.Prepare("SELECT role FROM users WHERE login=$1")
	if err != nil {
		tr.Rollback()
		return "", err
	}
	defer statement.Close()

	var role string
	err = statement.QueryRow(login).Scan(&role)
	if err != nil {
		tr.Rollback()
		if err == sql.ErrNoRows {
			return "", errors.ErrUserNotFound
		}
		return "", err
	}

	if err := tr.Commit(); err != nil {
		return "", err
	}

	return role, nil
}

func (u *Users) UserExists(ctx context.Context, login, email string) (bool, error) {
	tr, err := u.db.Begin()
	if err != nil {
		return false, err
	}
	statement, err := tr.Prepare("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 OR email = $2)")
	if err != nil {
		tr.Rollback()
		return false, err
	}
	defer statement.Close()

	var exists bool // Changed from string to bool
	err = statement.QueryRow(login, email).Scan(&exists)
	if err != nil {
		tr.Rollback()
		return false, err
	}

	if err := tr.Commit(); err != nil {
		return false, err
	}

	return exists, nil
}

// GetByLogin retrieves a user by login
func (u *Users) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	tr, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	
	statement, err := tr.Prepare("SELECT login, name, email, password, role FROM users WHERE login = $1")
	if err != nil {
		tr.Rollback()
		return nil, err
	}
	defer statement.Close()

	var user domain.User
	err = statement.QueryRow(login).Scan(&user.Login, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		tr.Rollback()
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	if err := tr.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}
