package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/service/exchangeRate"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("config init failed", zap.Error(err))
		os.Exit(1)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("tg client init failed", zap.Error(err))
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN())
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer conn.Close(ctx)

	expenseRepo := database.NewExpenseDb(conn)
	exchangeRateRepo := database.NewRateDb(conn)
	userRepo := database.NewUserDb(conn)
	txManager := database.NewTxManager(conn)

	exchangeRateService := exchangeRate.New(
		exchangeRateRepo,
		cfg.ExchangeRateAPIKey(),
		cfg.ExchangeRateBaseURI(),
		cfg.ExchangeRateRefreshRateInMin(),
	)
	exchangeRateService.Run()

	msgModel := messages.New(ctx, tgClient, expenseRepo, exchangeRateRepo, userRepo, txManager)
	go tgClient.ListenUpdates(msgModel)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	exchangeRateService.Close()
	tgClient.Close()
}
