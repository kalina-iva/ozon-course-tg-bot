package exchangeRate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
)

const timeoutInMin = 10

type currencyRepository interface {
	SaveRate(code string, rate float64)
}

type Service struct {
	currencyRepository currencyRepository
	cAPIKey            string
	baseURI            string
	refreshRateInMin   time.Duration
	wg                 sync.WaitGroup
	cancel             context.CancelFunc
}

type currencyResult struct {
	Rates map[string]float64 `json:"rates"`
}

func New(currencyRepo currencyRepository, cAPIKey string, baseURI string, refreshRateInMin int64) *Service {
	return &Service{
		currencyRepository: currencyRepo,
		cAPIKey:            cAPIKey,
		baseURI:            baseURI,
		refreshRateInMin:   time.Duration(refreshRateInMin) * time.Minute,
	}
}

func (s *Service) Run() {
	chanForResp := make(chan currencyResult)
	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())

	s.wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(s.refreshRateInMin):
				err := s.getRates(chanForResp)
				if err != nil {
					log.Println("cannot get rates: ", err)
				}
			case <-ctx.Done():
				close(chanForResp)
				s.wg.Done()
				return
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		for result := range chanForResp {
			for code, rate := range result.Rates {
				log.Printf("%s %.2f", code, rate)
				s.currencyRepository.SaveRate(code, rate)
			}
		}
		s.wg.Done()
	}()
}

func (s *Service) getRates(ch chan<- currencyResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutInMin*time.Minute)
	defer cancel()

	url := fmt.Sprintf(
		"%s?base=%s&symbols=%s",
		s.baseURI,
		messages.DefaultCurrencyCode,
		strings.Join(messages.AvailableCurrencies, ","),
	)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	request.Header.Set("apikey", s.cAPIKey)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "cannot execute request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Wrap(err, fmt.Sprintf("unexpected status code %d", res.StatusCode))
	}

	var result currencyResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return errors.Wrap(err, "decode response body")
	}
	log.Println("exchange rate was successfully received")
	ch <- result
	return nil
}

func (s *Service) Close() {
	s.cancel()
	s.wg.Wait()
}
