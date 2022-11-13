package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/api"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/consumer"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn, err := grpc.Dial(cfg.ReportServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("did not connect", zap.Error(err))
	}
	defer conn.Close()

	reportClient := api.NewReportClient(conn)

	databaseConn, err := pgx.Connect(ctx, cfg.DatabaseDSN())
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer conn.Close()

	expenseRepo := database.NewExpenseDb(databaseConn)
	exchangeRateRepo := database.NewRateDb(databaseConn)
	txManager := database.NewTxManager(databaseConn)

	generator := report.NewGenerator(expenseRepo, exchangeRateRepo, txManager)

	go func() {
		err := consumer.NewConsumerGroup(
			ctx,
			cfg.Kafka().BrokerList,
			cfg.Kafka().Report.ConsumerGroup,
			cfg.Kafka().Report.Topic,
			reportClient,
			generator,
		)
		if err != nil {
			logger.Fatal("cannot init consumer group", zap.Error(err))
		}
	}()
	logger.Info("consumer is starting")
	defer consumer.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
