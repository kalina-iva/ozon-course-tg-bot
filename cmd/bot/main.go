package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/currency"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/repository/memory"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/service/exchangeRate"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	repo := memory.New()
	currencyRepo := currency.New()

	exchangeRateService := exchangeRate.New(
		currencyRepo,
		cfg.ExchangeRateAPIKey(),
		cfg.ExchangeRateBaseURI(),
		cfg.ExchangeRateRefreshRateInMin(),
	)
	exchangeRateService.Run()

	msgModel := messages.New(tgClient, repo, currencyRepo)
	go tgClient.ListenUpdates(msgModel)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	exchangeRateService.Close()
	tgClient.Close()
}
