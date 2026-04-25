package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	action := flag.String("action", "up", "Migration action: up or down")
	flag.Parse()

	// Исправленный путь к .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("⚠️ Warning loading .env: %v", err)
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	fmt.Printf("🔍 POSTGRES_URL: %s\n", postgresURL)
	if postgresURL == "" {
		log.Fatal("❌ POSTGRES_URL is empty! Check .env file")
	}

	conn, err := pgx.Connect(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("❌ DB connect: %v", err)
	}
	defer conn.Close(context.Background())

	migrationsDir := "internal/db/migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("❌ Read migrations: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	if *action == "up" {
		for _, fname := range sqlFiles {
			path := filepath.Join(migrationsDir, fname)
			content, _ := os.ReadFile(path)
			_, err := conn.Exec(context.Background(), string(content))
			if err != nil {
				log.Fatalf("❌ Migration %s failed: %v", fname, err)
			}
			fmt.Printf("✅ Applied: %s\n", fname)
		}
		fmt.Println("🎉 All migrations applied successfully!")
	}
}
