package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/service/exchangeRate"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logging"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/tracing"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	err := logging.InitLogger()
	if err != nil {
		log.Fatal("logger init failed:", err)
	}
	defer logging.Close()

	cfg, err := config.New()
	if err != nil {
		zap.L().Fatal("config init failed", zap.Error(err))
	}

	err = tracing.InitTracing(cfg.ServiceName(), cfg.SamplingRatio())
	if err != nil {
		zap.L().Fatal("tracing init failed", zap.Error(err))
	}
	defer tracing.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8088", nil); err != nil {
			zap.L().Fatal("cannot start server for metrics", zap.Error(err))
		}
	}()

	tgClient, err := tg.New(cfg)
	if err != nil {
		zap.L().Fatal("tg client init failed", zap.Error(err))
	}

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN())
	if err != nil {
		zap.L().Fatal("cannot connect to database", zap.Error(err))
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
	defer exchangeRateService.Close()

	msgModel := messages.New(tgClient, expenseRepo, exchangeRateRepo, userRepo, txManager)
	go tgClient.ListenUpdates(ctx, msgModel)
	defer tgClient.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
