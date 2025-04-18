package main

import (
	"database/sql"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/jamoowen/quoteapi/internal/auth"
	"github.com/jamoowen/quoteapi/internal/cache"
	"github.com/jamoowen/quoteapi/internal/email"
	"github.com/jamoowen/quoteapi/internal/http"
	"github.com/jamoowen/quoteapi/internal/jobs"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load .env", err)
	}
}

func main() {
	logger := log.Default()

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Startup error: %v", err.Error())
	}

	csvPath := path.Join(wd, "/data/quotes.csv")
	quoteCache, err := cache.NewQuoteCache(csvPath)
	if err != nil {
		logger.Fatal(err)
	}

	emailAddress := os.Getenv("SMTP_EMAIL_ADDRESS")
	emailPassword := os.Getenv("GMAIL_APP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpService, err := email.NewSmtpMailer(emailAddress, emailPassword, smtpHost, smtpPort)
	if err != nil {
		logger.Fatal(err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./db/quotedb.sqlite"
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(3)
	db.SetConnMaxIdleTime(5 * time.Minute)

	otpSecondsValid, err := strconv.Atoi(os.Getenv("OTP_SECONDS_VALID"))
	if err != nil {
		log.Fatal(err)
	}

	authSecret := os.Getenv("AUTH_SECRET")
	if authSecret == "" {
		log.Fatal("AUTH_SECRET env var missing")
	}
	authService := auth.NewAuthService(db, int64(otpSecondsValid), authSecret)

	requiredSecondIntervalsBeforeIpRateLimit := int64(5)
	requiredSecondIntervalsBeforeApiKeyRateLimit := int64(5)

	s, err := http.NewServer(quoteCache, smtpService, authService, requiredSecondIntervalsBeforeIpRateLimit, requiredSecondIntervalsBeforeApiKeyRateLimit, logger)
	if err != nil {
		log.Fatal(err)
	}
	cleanupInterval := 12 * time.Hour
	go jobs.CacheCleanupJob(cleanupInterval, int64(otpSecondsValid), authService, s)

	log.Fatal(s.StartServer())
}
