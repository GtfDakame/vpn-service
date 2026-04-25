package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateUser создаёт нового пользователя и сразу генерирует ему VPN UUID
func CreateUser(ctx context.Context, conn *pgx.Conn, tgID *int64) (*User, error) {
	vpnUUID := uuid.New()
	query := `
		INSERT INTO users (id, telegram_id, vpn_uuid, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, telegram_id, email, vpn_uuid, vpn_config_url, vpn_last_used_at, created_at, updated_at
	`

	var user User
	err := conn.QueryRow(ctx, query, uuid.New(), tgID, vpnUUID).Scan(
		&user.ID, &user.TelegramID, &user.Email, &user.VPNUUID,
		&user.VPNConfigURL, &user.VPNLastUsedAt, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

// GetUserByTelegramID ищет пользователя по Telegram ID
func GetUserByTelegramID(ctx context.Context, conn *pgx.Conn, tgID int64) (*User, error) {
	query := `
		SELECT id, telegram_id, email, vpn_uuid, vpn_config_url, vpn_last_used_at, created_at, updated_at
		FROM users WHERE telegram_id = $1
	`

	var user User
	err := conn.QueryRow(ctx, query, tgID).Scan(
		&user.ID, &user.TelegramID, &user.Email, &user.VPNUUID,
		&user.VPNConfigURL, &user.VPNLastUsedAt, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

// UpdateUserLastUsed обновляет время последнего использования
func UpdateUserLastUsed(ctx context.Context, conn *pgx.Conn, userID uuid.UUID) error {
	_, err := conn.Exec(ctx,
		`UPDATE users SET vpn_last_used_at = NOW(), updated_at = NOW() WHERE id = $1`,
		userID)
	return err
}
