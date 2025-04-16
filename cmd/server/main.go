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
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
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

	db, err := sql.Open("sqlite3", "./db/quotedb.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	otpSecondsValid, err := strconv.Atoi(os.Getenv("OTP_SECONDS_VALID"))
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewAuthService(db, int64(otpSecondsValid))

	requiredSecondIntervalsBeforeIpRateLimit := int64(5)
	requiredSecondIntervalsBeforeApiKeyRateLimit := int64(5)
	s, err := http.NewServer(quoteCache, smtpService, authService, requiredSecondIntervalsBeforeIpRateLimit, requiredSecondIntervalsBeforeApiKeyRateLimit, logger)
	if err != nil {
		log.Fatal(err)
	}

	authCacheCleanupIntervals := time.Duration(5 * time.Minute)
	go http.CleanupAuthCache(&s.ApiKeyRateLimiter, &s.IpRateLimiter, authCacheCleanupIntervals, requiredSecondIntervalsBeforeApiKeyRateLimit, requiredSecondIntervalsBeforeIpRateLimit)

	otpCacheCleanupIntervals := time.Duration(5 * time.Minute)
	go cache.CleanupOtpCache()

	log.Fatal(s.StartServer())
}
