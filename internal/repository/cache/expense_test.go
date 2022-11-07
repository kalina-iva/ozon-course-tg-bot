package cache

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	cacheMocks "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/mocks/cache"
	messagesMocks "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/mocks/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

func Test_OnCreate_ShouldInvalidateTags(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")

	expenseRepo.EXPECT().Create(gomock.Any(), int64(123), "category_name", uint64(8700), mockTime).Return(nil)
	cacheManager.EXPECT().Invalidate(gomock.Any(), []string{"report123"})

	err := cache.Create(context.Background(), int64(123), "category_name", uint64(8700), mockTime)

	assert.NoError(t, err)
}

func Test_OnCreate_CannotCreateInDB(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	mockErr := errors.New("new error")

	expenseRepo.EXPECT().Create(gomock.Any(), int64(123), "category_name", uint64(8700), mockTime).Return(mockErr)
	cacheManager.EXPECT().Invalidate(gomock.Any(), []string{"report123"}).Times(0)

	err := cache.Create(context.Background(), int64(123), "category_name", uint64(8700), mockTime)
	assert.Error(t, err)
}

func Test_OnReport_FromCache(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	expiry := 24 * time.Hour

	cacheManager.EXPECT().GetBytes(gomock.Any(), "report123:07/12/2012").
		Times(1).
		Return([]byte("[{\"Category\":\"new cat\",\"AmountInKopecks\":8700}]"), nil)

	expenseRepo.EXPECT().Report(gomock.Any(), int64(123), mockTime).Times(0)
	cacheManager.EXPECT().Set(gomock.Any(), "report123:07/12/2012", nil, []string{"report123"}, expiry).Times(0)

	_, err := cache.Report(context.Background(), int64(123), mockTime)
	assert.NoError(t, err)
}

func Test_OnReport_FromDB(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	expiry := 24 * time.Hour
	report := []*entity.Report{
		{
			Category:        "cat",
			AmountInKopecks: 8700,
		},
	}

	cacheManager.EXPECT().GetBytes(gomock.Any(), "report123:07/12/2012").Times(1).Return(nil, errNotFoundInCache)

	expenseRepo.EXPECT().Report(gomock.Any(), int64(123), mockTime).Times(1).Return(report, nil)
	cacheManager.EXPECT().
		Set(gomock.Any(), "report123:07/12/2012", []byte("[{\"Category\":\"cat\",\"AmountInKopecks\":8700}]"), []string{"report123"}, expiry).
		Times(1)

	_, err := cache.Report(context.Background(), int64(123), mockTime)
	assert.NoError(t, err)
}

func Test_OnReport_ErrFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	expiry := 24 * time.Hour
	mockErr := errors.New("new error")

	cacheManager.EXPECT().GetBytes(gomock.Any(), "report123:07/12/2012").Times(1).Return(nil, errNotFoundInCache)

	expenseRepo.EXPECT().Report(gomock.Any(), int64(123), mockTime).Times(1).Return(nil, mockErr)
	cacheManager.EXPECT().Set(gomock.Any(), "report123:07/12/2012", nil, []string{"report123"}, expiry).Times(0)

	_, err := cache.Report(context.Background(), int64(123), mockTime)
	assert.Error(t, err)
}

func Test_GetAmountByPeriod_OnOk(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")

	expenseRepo.EXPECT().GetAmountByPeriod(gomock.Any(), int64(123), mockTime).Times(1)

	_, err := cache.GetAmountByPeriod(context.Background(), int64(123), mockTime)
	assert.NoError(t, err)
}

func Test_GetAmountByPeriod_ErrFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)

	expenseRepo := messagesMocks.NewMockExpenseRepository(ctrl)
	cacheManager := cacheMocks.NewMockcacheManager(ctrl)
	cache := NewExpenseCache(expenseRepo, cacheManager)

	mockTime, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	mockErr := errors.New("new error")

	expenseRepo.EXPECT().GetAmountByPeriod(gomock.Any(), int64(123), mockTime).Times(1).Return(uint64(0), mockErr)

	_, err := cache.GetAmountByPeriod(context.Background(), int64(123), mockTime)
	assert.Error(t, err)
}
