package service

import (
	"log/slog"
)

type UserService struct {
	logger *slog.Logger
}

func NewUserService() *UserService {
	return &UserService{
		logger: slog.Default(),
	}
}

func (s *UserService) GetUserByTelegramID(telegramID int64) (*User, error) {
	// TODO: Implement with sqlc
	s.logger.Info("Getting user", "telegram_id", telegramID)
	return nil, nil
}

func (s *UserService) CreateUser(telegramID int64, username string, firstName string, lastName string) (*User, error) {
	// TODO: Implement with sqlc
	s.logger.Info("Creating user", "telegram_id", telegramID, "username", username)
	return nil, nil
}

type User struct {
	ID        int64
	TelegramID int64
	Username   string
	FirstName  string
	LastName   string
	CreatedAt  string
	UpdatedAt  string
}
