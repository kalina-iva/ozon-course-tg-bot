package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

const expiration = 24 * time.Hour

type ExpenseCache struct {
	db      messages.ExpenseRepository
	manager *Manager
}

func NewExpenseCache(db messages.ExpenseRepository, manager *Manager) *ExpenseCache {
	return &ExpenseCache{
		db:      db,
		manager: manager,
	}
}

func (c *ExpenseCache) New(ctx context.Context, userID int64, category string, amount uint64, date time.Time) error {
	if err := c.db.New(ctx, userID, category, amount, date); err != nil {
		return errors.Wrap(err, "cannot create expense in db")
	}
	tag := getReportTag(userID)
	c.manager.Invalidate(ctx, []string{tag})
	return nil
}

func (c *ExpenseCache) Report(ctx context.Context, userID int64, period time.Time) ([]*entity.Report, error) {
	key := getReportKey(userID, period)
	reportData, err := c.reportFromCache(ctx, key)
	if err == nil {
		return reportData, nil
	}
	if !errors.Is(err, errNotFoundInCache) {
		logger.Error("cannot get report data from cache", zap.Error(err))
	}

	reportData, err = c.db.Report(ctx, userID, period)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get report data from db")
	}

	tag := getReportTag(userID)
	if err = c.saveReportToCache(ctx, reportData, key, tag); err != nil {
		logger.Error("cannot save report to cache", zap.Error(err))
	}

	return reportData, nil
}

func getReportKey(userID int64, period time.Time) string {
	var sb strings.Builder
	sb.WriteString("report")
	sb.WriteString(strconv.FormatInt(userID, 10))
	sb.WriteString(":")
	sb.WriteString(period.Format("01/02/2006"))
	return sb.String()
}

func getReportTag(userID int64) string {
	var sb strings.Builder
	sb.WriteString("report")
	sb.WriteString(strconv.FormatInt(userID, 10))
	return sb.String()
}

func (c *ExpenseCache) reportFromCache(ctx context.Context, key string) ([]*entity.Report, error) {
	data, err := c.manager.GetBytes(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get data from cache manager")
	}

	var reportData []*entity.Report
	if err = json.Unmarshal(data, &reportData); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal report data")
	}
	return reportData, nil
}

func (c *ExpenseCache) saveReportToCache(ctx context.Context, reportData []*entity.Report, key string, tag string) error {
	marshaledData, err := json.Marshal(reportData)
	if err != nil {
		return errors.Wrap(err, "cannot marshal report")
	}
	if err = c.manager.Set(ctx, key, marshaledData, []string{tag}, expiration); err != nil {
		return errors.Wrap(err, "cannot save report to cache")
	}
	return nil
}

func (c *ExpenseCache) GetAmountByPeriod(ctx context.Context, userID int64, period time.Time) (uint64, error) {
	amount, err := c.db.GetAmountByPeriod(ctx, userID, period)
	if err != nil {
		return amount, errors.Wrap(err, "cannot get amount from db")
	}
	return amount, nil
}
