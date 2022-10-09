package exchangeRate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	currencyAPIKey     string
}

type currencyResult struct {
	Rates map[string]float64 `json:"rates"`
}

func New(currencyRepo currencyRepository, currencyAPIKey string) *Service {
	return &Service{
		currencyRepository: currencyRepo,
		currencyAPIKey:     currencyAPIKey,
	}
}

func (s *Service) Run() {
	chanForResp := make(chan currencyResult)
	go func() {
		for {
			err := s.getRates(chanForResp)
			if err != nil {
				log.Println("cannot get rates: ", err)
			}
			time.Sleep(time.Minute)
		}
	}()

	go func() {
		for result := range chanForResp {
			for code, rate := range result.Rates {
				log.Printf("%s %.2f", code, rate)
				s.currencyRepository.SaveRate(code, rate)
			}
		}
	}()
}

func (s *Service) getRates(ch chan<- currencyResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutInMin*time.Minute)
	defer cancel()

	url := fmt.Sprintf(
		"https://api.apilayer.com/fixer/latest?base=%s&symbols=%s",
		messages.DefaultCurrencyCode,
		strings.Join(messages.AvailableCurrencies, ","),
	)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}
	request.Header.Set("apikey", s.currencyAPIKey)

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
