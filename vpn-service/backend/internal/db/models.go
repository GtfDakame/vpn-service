package db

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя в БД
type User struct {
	ID            uuid.UUID  `json:"id"`
	TelegramID    *int64     `json:"telegram_id,omitempty"`
	Email         *string    `json:"email,omitempty"`
	VPNUUID       *uuid.UUID `json:"vpn_uuid,omitempty"`
	VPNConfigURL  *string    `json:"vpn_config_url,omitempty"`
	VPNLastUsedAt *time.Time `json:"vpn_last_used_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Status        *string    `json:"status,omitempty"`
}
