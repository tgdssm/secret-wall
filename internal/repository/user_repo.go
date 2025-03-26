package repository

import (
	"database/sql"
	"errors"
	"secretWall/internal/domain"
	"strings"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) FindByAppleSub(appleSub string) (*domain.User, error) {
	statement, err := s.db.Prepare("SELECT id, apple_sub, created_at FROM users WHERE apple_sub = $1")
	if err != nil {
		return nil, err
	}

	defer statement.Close()

	row := statement.QueryRow(appleSub)
	var user domain.User
	if err := row.Scan(&user.ID, &user.AppleSub, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserRepo) Create(user *domain.User) error {
	statement, err := s.db.Prepare("INSERT INTO users(id, apple_sub, created_at) VALUES ($1, $2, $3)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, user.AppleSub, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrUserAlreadyExist
		}
		return err
	}

	return nil
}
