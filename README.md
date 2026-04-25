# 🌐 HyperInet VPN Service

Полноценный сервис управления VPN-подписками с автоматической выдачей ключей, биллингом и защитой от шаринга.

## 🛠️ Стек
- **Backend**: Go 1.21, Fiber, pgx, telebot
- **VPN**: Xray-core (VLESS + Reality)
- **Database**: PostgreSQL
- **Frontend**: React + Tailwind (mini-app)
- **Payments**: Platega.io / Cryptomus

## ✨ Функционал
- 🔑 Автоматическая генерация VLESS+Reality ключей
- 💳 Пополнение баланса (карты, СБП, крипта)
- 📱 Telegram-бот с полным циклом: покупка, управление, поддержка
- 🌐 Веб-кабинет: авторизация, ключи, баланс, привязка аккаунтов
- 🛡️ Анти-шаринг: блокировка при подключении с 2+ IP
- 🔄 Авто-продление: отключение/восстановление ключей при 0 ₽

## 🚀 Быстрый запуск (Dev)
```bash
# 1. Клон и зависимости
git clone https://github.com/GtfDakame/vpn-service.git
cd vpn-service
go mod tidy

# 2. Настройка окружения
cp infra/.env.example infra/.env  # заполни своими данными

# 3. Запуск
go run backend/cmd/api/main.go
