# 🔍 Смотреть логи бота в реальном времени
sudo journalctl -u vpn-bot -f

# 🔄 Перезапустить бота
sudo systemctl restart vpn-bot

# 🛑 Остановить бота
sudo systemctl stop vpn-bot

# ✅ Проверить статус
sudo systemctl status vpn-bot --no-pager

# 🐳 Проверить контейнеры БД
docker compose -f /root/vpn-service/infra/docker-compose.yml ps

# 👥 Посмотреть список пользователей в Xray
cat /opt/xray/uuids.txt

# 🗄️ Применить миграции вручную
cd /root/vpn-service/backend && go run cmd/migrate/main.go -action=up

# 🧹 Очистить старые зависимости (если снова ошибка go.mod)
cd /root/vpn-service/backend && go mod tidy && sudo systemctl restart vpn-bot

# Быстрая очистка uuids.txt от повторов
sort -u /opt/xray/uuids.txt -o /opt/xray/uuids.txt && /opt/xray/sync_users.sh