package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	"vpn-service/internal/db"
	"vpn-service/internal/vpn"

	telebot "gopkg.in/tucnak/telebot.v2"

	"github.com/jackc/pgx/v5"
)

type Bot struct {
	*telebot.Bot
	conn *pgx.Conn
}

func NewBot(token string, conn *pgx.Conn) (*Bot, error) {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}
	return &Bot{Bot: b, conn: conn}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	// === /start ===
	b.Handle("/start", func(m *telebot.Message) {
		tgID := m.Sender.ID
		log.Printf("👤 /start from TG: %d", tgID)

		user, err := db.GetUserByTelegramID(ctx, b.conn, tgID)
		if err != nil {
			// Создаём нового пользователя (UUID генерируется автоматически в DB)
			user, err = db.CreateUser(ctx, b.conn, &tgID)
			if err != nil {
				b.Send(m.Sender, "❌ Ошибка регистрации. Попробуйте позже.")
				log.Printf("❌ Failed to create user: %v", err)
				return
			}
			b.Send(m.Sender, fmt.Sprintf(
				"👋 Привет! Вы зарегистрированы.\n"+
					"Ваш ID: `%s`\n\n"+
					"Отправьте /config чтобы получить ссылку для HAPP.",
				user.ID.String(),
			), telebot.ModeMarkdown)
		} else {
			b.Send(m.Sender, fmt.Sprintf(
				"👋 С возвращением!\n"+
					"Ваш ID: `%s`\n\n"+
					"Отправьте /config чтобы получить ссылку.",
				user.ID.String(),
			), telebot.ModeMarkdown)
		}
	})

	// === /config ===
	b.Handle("/config", func(m *telebot.Message) {
		tgID := m.Sender.ID
		log.Printf("🔗 /config from TG: %d", tgID)

		user, err := db.GetUserByTelegramID(ctx, b.conn, tgID)
		if err != nil {
			b.Send(m.Sender, "❌ Сначала отправьте /start", telebot.ModeMarkdown)
			return
		}

		// На всякий случай проверяем, есть ли VPN UUID
		if user.VPNUUID == nil {
			b.Send(m.Sender, "⏳ Конфиг ещё не сгенерирован. Отправьте /start заново.", telebot.ModeMarkdown)
			return
		}

		serverHost := "83.97.78.70"
		serverPort := 443

		fixedUUID := "167ea650-60cc-43e9-91b9-0a1b978125fd" // ← тот же, что в Xray config.json
		vlessLink := vpn.GenerateVLESSURI(
			fixedUUID,
			serverHost,
			serverPort,
		)

		_, err = b.Send(m.Sender, fmt.Sprintf(
			"🔗 Ваша ссылка для HAPP:\n\n"+
				"```\n%s\n```\n\n"+
				"✅ Инструкция:\n"+
				"1. Скопируйте ссылку целиком\n"+
				"2. Откройте HAPP → + → Импорт из буфера\n"+
				"3. Включите профиль\n\n"+
				"⚠️ Убедитесь, что порт 443 не заблокирован.",
			vlessLink,
		), telebot.ModeMarkdown)

		if err != nil {
			log.Printf("❌ Failed to send config: %v", err)
		}

		_ = db.UpdateUserLastUsed(ctx, b.conn, user.ID)
	})

	// === /help ===
	b.Handle("/help", func(m *telebot.Message) {
		b.Send(m.Sender,
			"📋 Доступные команды:\n"+
				"/start — Регистрация или вход\n"+
				"/config — Получить ссылку для HAPP\n"+
				"/help — Эта справка",
			telebot.ModeMarkdown,
		)
	})

	log.Println("🤖 Telegram bot starting...")
	b.Bot.Start()
	return nil
}

func (b *Bot) Stop() {
	if b.Bot != nil {
		b.Bot.Stop()
	}
}
