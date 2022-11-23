package report

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/mocks/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

func Test_OnGenerateReport_onOk(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := mocks.NewMockexpenseRepository(ctrl)
	rateRepository := mocks.NewMockexchangeRateRepository(ctrl)
	txManager := mocks.NewMocktxManager(ctrl)
	model := NewGenerator(expenseRepo, rateRepository, txManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")

	rateRepository.EXPECT().GetRate(gomock.Any(), "RUB").Return(1.1, nil)
	expenseRepo.EXPECT().Report(gomock.Any(), int64(123), mockTime).Return([]*entity.Report{
		{
			Category:        "cat",
			AmountInKopecks: 10000,
		},
	}, nil)

	report := model.GenerateReport(context.Background(), Request{
		UserID:       123,
		Period:       mockTime,
		CurrencyCode: "RUB",
	})
	assert.Equal(t, "cat: 110.00 ₽\n", report)
}
