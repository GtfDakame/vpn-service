package main

import (
	"context"
	"log"
	"os"

	"vpn-service/internal/db"
	"vpn-service/internal/telegram"
	"vpn-service/internal/vpn"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("../.env")

	// 1. Подключение к PostgreSQL
	connString := os.Getenv("POSTGRES_URL")
	if connString == "" {
		// Fallback с правильными данными из .env
		connString = "postgres://vpnadmin:dev_pass_2026!@localhost:5432/vpn_service?sslmode=disable"
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)
	log.Println("✅ PostgreSQL connected successfully")

	// 2. Инициализация Fiber
	app := fiber.New()
	app.Use(cors.New())

	// 3. API: Генерация конфига по Telegram ID
	app.Post("/api/config", func(c *fiber.Ctx) error {
		type Request struct {
			TelegramID int64 `json:"telegram_id"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Ищем или создаём пользователя
		user, err := db.GetUserByTelegramID(ctx, conn, req.TelegramID)
		if err != nil {
			user, err = db.CreateUser(ctx, conn, &req.TelegramID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
			}
		}

		// Генерируем ссылку (3 аргумента: uuid, host, port)
		link := vpn.GenerateVLESSURI(user.VPNUUID.String(), "83.97.78.70", 443)

		return c.JSON(fiber.Map{
			"success": true,
			"config":  link,
			"user_id": user.ID.String(),
		})
	})

	// 4. Запуск API сервера
	go func() {
		log.Println("🚀 API Server starting on :8080")
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("❌ API server failed: %v", err)
		}
	}()

	// 5. Запуск Telegram бота
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgToken == "" {
		log.Println("⚠️ TELEGRAM_BOT_TOKEN not set. Bot won't start.")
	} else {
		bot, err := telegram.NewBot(tgToken, conn)
		if err != nil {
			log.Fatalf("❌ Failed to create bot: %v", err)
		}

		log.Println("🤖 Starting Telegram bot...")
		if err := bot.Start(ctx); err != nil {
			log.Fatalf("❌ Bot failed: %v", err)
		}
	}
	
	// Блокируем основной поток
	select {}
}
