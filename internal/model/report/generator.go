package report

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

const cntKopInRub = 100

type exchangeRateRepository interface {
	GetRate(ctx context.Context, code string) (float64, error)
}

type expenseRepository interface {
	Report(ctx context.Context, userID int64, period time.Time) ([]*entity.Report, error)
}

type txManager interface {
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type Request struct {
	UserID       int64     `json:"user_id"`
	Period       time.Time `json:"period"`
	CurrencyCode string    `json:"currency_code"`
}

type Generator struct {
	expenseRepo      expenseRepository
	exchangeRateRepo exchangeRateRepository
	txManager        txManager
}

func NewGenerator(
	expenseRepo expenseRepository,
	exchangeRateRepo exchangeRateRepository,
	txManager txManager,
) *Generator {
	return &Generator{
		expenseRepo:      expenseRepo,
		exchangeRateRepo: exchangeRateRepo,
		txManager:        txManager,
	}
}

func (r *Generator) GenerateReport(ctx context.Context, request Request) string {
	rate, err := r.exchangeRateRepo.GetRate(ctx, request.CurrencyCode)
	if err != nil {
		logger.Error("cannot get rate from expenseRepo", zap.Error(err))
		return canNotGetRate
	}

	currencyShort := getCurrencyShortByCode(request.CurrencyCode)
	report, err := r.expenseRepo.Report(ctx, request.UserID, request.Period)
	if err != nil {
		logger.Error("cannot get report request", zap.Error(err))
		return canNotCreateReport
	}

	var sb strings.Builder
	for _, item := range report {
		amount := float64(item.AmountInKopecks) * rate
		sb.WriteString(item.Category)
		sb.WriteString(fmt.Sprintf(": %.2f %v\n", amount/cntKopInRub, currencyShort))
	}
	if sb.Len() == 0 {
		logger.Info("no request for report")
		return noDataForReport
	}
	return sb.String()
}

func getCurrencyShortByCode(code string) (short string) {
	switch code {
	case "RUB":
		short = "₽"
	case "USD":
		short = "＄"
	case "EUR":
		short = "€"
	case "CNY":
		short = "元"
	}
	return
}
