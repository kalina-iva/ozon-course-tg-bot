package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/cache"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/service/exchangeRate"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/metrics"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/tracing"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	err = logger.InitLogger(cfg.ServiceEnv())
	if err != nil {
		log.Fatal("logger init failed:", err)
	}
	defer logger.Close()

	err = tracing.InitTracing(cfg.ServiceName(), cfg.SamplingRatio())
	if err != nil {
		logger.Fatal("tracing init failed", zap.Error(err))
	}
	defer tracing.Close()

	go func() {
		err = metrics.InitMetrics(cfg.MetricsServerAddress())
		if err != nil {
			logger.Fatal("metrics init failed", zap.Error(err))
		}
	}()
	defer metrics.Close()

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN())
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer conn.Close(ctx)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost(),
	})

	expenseRepo := database.NewExpenseDb(conn)
	exchangeRateRepo := database.NewRateDb(conn)
	userRepo := database.NewUserDb(conn)
	txManager := database.NewTxManager(conn)

	cacheManager := cache.NewManager(redisClient)
	ExpenseCache := cache.NewExpenseCache(expenseRepo, cacheManager)

	exchangeRateService := exchangeRate.New(
		exchangeRateRepo,
		cfg.ExchangeRateAPIKey(),
		cfg.ExchangeRateBaseURI(),
		cfg.ExchangeRateRefreshRateInMin(),
	)
	exchangeRateService.Run()
	defer exchangeRateService.Close()

	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatal("tg client init failed", zap.Error(err))
	}

	msgModel := messages.New(tgClient, ExpenseCache, exchangeRateRepo, userRepo, txManager)
	go tgClient.ListenUpdates(ctx, msgModel)
	defer tgClient.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
