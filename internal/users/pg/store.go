package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/nskondratev/go-telegram-translator-bot/internal/users"
)

const (
	selectUserByTelegramIDQuery = `
SELECT
	telegram_user_id,
	username,
	first_name,
	last_name,
	language,
	source_lang,
	target_lang,
	points
FROM "users"
WHERE telegram_user_id = $1
`
	insertUserQuery = `
INSERT INTO "users"
	(telegram_user_id, username, first_name, last_name, language, source_lang, target_lang, points)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
`
)

type Store struct {
	db *pgx.ConnPool
}

func New(db *pgx.ConnPool) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByTelegramUserID(ctx context.Context, tgUserID int64) (users.User, error) {
	var user users.User
	row := s.db.QueryRowEx(ctx, selectUserByTelegramIDQuery, &pgx.QueryExOptions{}, tgUserID)
	err := row.Scan(&user.TelegramUserID, &user.UserName, &user.FirstName, &user.LastName, &user.Lang, &user.SourceLang,
		&user.TargetLang, &user.Points)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, users.ErrUserNotFound
		}
		return user, fmt.Errorf("failed to scan user by telegram user id: %w", err)
	}
	return user, nil
}

func (s *Store) StoreUser(ctx context.Context, user *users.User) error {
	_, err := s.db.ExecEx(ctx, insertUserQuery, &pgx.QueryExOptions{},
		user.TelegramUserID, user.UserName, user.FirstName, user.LastName, user.Lang, user.SourceLang, user.TargetLang,
		user.Points)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}