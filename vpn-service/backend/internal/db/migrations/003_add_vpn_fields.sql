-- Migration: 003_add_vpn_fields.sql
-- Description: Добавляет поля для хранения VPN-конфигов пользователей

-- Добавляем колонки в таблицу users
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS vpn_uuid UUID DEFAULT gen_random_uuid() UNIQUE,
ADD COLUMN IF NOT EXISTS vpn_config_url TEXT,
ADD COLUMN IF NOT EXISTS vpn_last_used_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Создаём индекс для быстрого поиска по vpn_uuid
CREATE INDEX IF NOT EXISTS idx_users_vpn_uuid ON users(vpn_uuid);

-- Создаём таблицу для статистики трафика (опционально, для будущего)
CREATE TABLE IF NOT EXISTS user_traffic (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    upload_bytes BIGINT DEFAULT 0,
    download_bytes BIGINT DEFAULT 0,
    session_start TIMESTAMP DEFAULT NOW(),
    session_end TIMESTAMP,
    server_ip INET,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для таблицы трафика
CREATE INDEX IF NOT EXISTS idx_user_traffic_user_id ON user_traffic(user_id);
CREATE INDEX IF NOT EXISTS idx_user_traffic_session ON user_traffic(session_start, session_end);

-- Комментарий к миграции
COMMENT ON TABLE user_traffic IS 'Статистика использования VPN пользователями';