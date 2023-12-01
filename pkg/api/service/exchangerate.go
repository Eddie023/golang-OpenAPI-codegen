package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/shopspring/decimal"
)

type ExchangeRateGetter struct{}

type ExchangeRatePayload struct {
	CountryName string
	Currency    string
	RecordDate  time.Time
}

type ExchangeRateResponse struct {
	CountryCurrencyDesc string `json:"country_currency_desc"`
	ExchangeRate        string `json:"exchange_rate"`
	RecordDate          string `json:"record_date"`
}

type ExchangeRateAPIResponse struct {
	Data []ExchangeRateResponse `json:"data"`
}

const TREASURY_RATES_OF_EXCHANGE_API_URL = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"

func (e *ExchangeRateGetter) GetExchangeRate(ctx context.Context, payload ExchangeRatePayload) (ExchangeRateResponse, error) {

	req, err := http.NewRequest("GET", TREASURY_RATES_OF_EXCHANGE_API_URL, nil)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	filter := fmt.Sprintf("record_date:lte:%s,country_currency_desc:eq:%s-%s", payload.RecordDate.Format(time.DateOnly), payload.CountryName, payload.Currency)
	fields := "country_currency_desc,exchange_rate,record_date"
	// sort by record_date in descending order such that we will get the first item which is closest to our purchase date within last six months
	sort := "-record_date"

	req.URL.RawQuery = fmt.Sprintf("filter=%s&fields=%s&sort=%s&page[size]=1", filter, fields, sort)

	client := &http.Client{}

	var resp *http.Response

	// since, this API have a rate limiting, we will try to expotentially backoff and retry if we get too many request error from the API.
	operation := func() error {
		resp, err = client.Do(req)
		if err != nil {
			return &backoff.PermanentError{}
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("too many request error")
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 1 * time.Minute

	err = backoff.Retry(operation, expBackoff)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Debug("exchange request failed", "status_code", resp.StatusCode)
		return ExchangeRateResponse{}, apiout.NewRequestError(fmt.Errorf("the exchange rate service failed with status code %v", resp.StatusCode), http.StatusInternalServerError)
	}

	var response ExchangeRateAPIResponse

	slog.Debug("generated exchange rate API", "url", req.URL)

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	// for invalid country or currency, API will still return 200 with empty list
	if len(response.Data) == 0 {
		return ExchangeRateResponse{}, apiout.NewRequestError(errors.New("the purchase cannot be converted to the target currency"), http.StatusBadRequest)
	}

	if resp != nil {
		resp.Body.Close()
	}

	// parse string to Date
	recordDate, err := time.Parse(time.DateOnly, response.Data[0].RecordDate)
	if err != nil {
		return ExchangeRateResponse{}, apiout.NewRequestError(errors.New("unable to parse returned record date"), http.StatusInternalServerError)
	}

	// currency conversion rate can be less than or equal to purchase date from within the last 6 months
	sixMonthBeforePurchaseDate := getSixMonthBeforePurchaseDate(payload.RecordDate)
	if recordDate.Before(sixMonthBeforePurchaseDate) {
		return ExchangeRateResponse{}, apiout.NewRequestError(errors.New("the purchase cannot be converted to the target currency"), http.StatusBadRequest)
	}

	// we can return the first item since we have already sorted our API response to our need.
	return response.Data[0], nil
}

func (e *ExchangeRateGetter) ConvertCurrency(payload ExchangeRatePayload, trans *ent.Transaction, er ExchangeRateResponse) (types.GetPurchaseTransaction, error) {
	exchangeRate, err := decimal.NewFromString(er.ExchangeRate)
	if err != nil {
		return types.GetPurchaseTransaction{}, err
	}

	convertedAmount := convertAmount(trans.AmountInUsd, exchangeRate)

	response := types.GetPurchaseTransaction{
		ConvertedPurchasePrice: struct {
			Amount           string "json:\"amount\""
			Country          string "json:\"country\""
			Currency         string "json:\"currency\""
			ExchangeRateDate string "json:\"exchangeRateDate\""
			ExchangeRateUsed string "json:\"exchangeRateUsed\""
		}{
			Amount:           RoundToNearestCent(convertedAmount).String(),
			Country:          payload.CountryName,
			Currency:         payload.Currency,
			ExchangeRateDate: er.RecordDate,
			ExchangeRateUsed: er.ExchangeRate,
		},
		Description: trans.Description,
		OriginalPurchasePrice: struct {
			Amount   string "json:\"amount\""
			Currency string "json:\"currency\""
		}{
			Amount:   trans.AmountInUsd.String(),
			Currency: "USD",
		},
		Transaction: struct {
			Date string "json:\"date\""
			Id   string "json:\"id\""
		}{
			Date: trans.Date.Format(time.DateTime),
			Id:   trans.ID.String(),
		},
	}

	return response, nil
}

func getSixMonthBeforePurchaseDate(d time.Time) time.Time {
	return d.AddDate(0, -6, 0)
}

func convertAmount(original decimal.Decimal, exchangeRate decimal.Decimal) decimal.Decimal {
	return original.Mul(exchangeRate)
}
