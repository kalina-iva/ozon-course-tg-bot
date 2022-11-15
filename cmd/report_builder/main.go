package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	protoReport "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/api/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/consumer"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/helper/grpcconn"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/cache"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/database"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn

func newReportClient(serverAddr string) (protoReport.ReportClient, error) {
	var err error
	conn, err = grpcconn.NewClientConn(serverAddr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create report client connection")
	}

	return protoReport.NewReportClient(conn), nil
}

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

	reportClient, err := newReportClient(cfg.ReportServerAddress())
	if err != nil {
		logger.Fatal("cannot create grpc report client", zap.Error(err))
	}
	defer conn.Close()

	databaseConn, err := pgx.Connect(ctx, cfg.DatabaseDSN())
	if err != nil {
		logger.Fatal("cannot connect to database", zap.Error(err))
	}
	defer databaseConn.Close(ctx)

	expenseRepo := database.NewExpenseDb(databaseConn)
	exchangeRateRepo := database.NewRateDb(databaseConn)
	txManager := database.NewTxManager(databaseConn)

	cacheManager := cache.NewManager(cfg.RedisHost())
	ExpenseCache := cache.NewExpenseCache(expenseRepo, cacheManager)

	generator := report.NewGenerator(ExpenseCache, exchangeRateRepo, txManager)

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
