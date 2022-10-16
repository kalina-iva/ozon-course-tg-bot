package main

import (
	"context"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/memory"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/service/exchangeRate"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN())
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer conn.Close(context.Background())
	//conn.Ping()

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	repo := memory.NewExpense()
	currencyRepo := memory.NewCurrency()
	userRepo := database.NewUserDb(conn)

	exchangeRateService := exchangeRate.New(
		currencyRepo,
		cfg.ExchangeRateAPIKey(),
		cfg.ExchangeRateBaseURI(),
		cfg.ExchangeRateRefreshRateInMin(),
	)
	exchangeRateService.Run()

	msgModel := messages.New(tgClient, repo, currencyRepo, userRepo)
	go tgClient.ListenUpdates(msgModel)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	exchangeRateService.Close()
	tgClient.Close()
}
