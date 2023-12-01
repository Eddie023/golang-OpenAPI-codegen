package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGetExchangeRate(t *testing.T) {
	tests := []struct {
		name         string
		purchaseDate string
		payload      ExchangeRatePayload
		want         ExchangeRateResponse
		wantErr      bool
	}{
		{
			name:         "should return correct exchange rate for valid country name and currency",
			purchaseDate: "2022-11-30",
			payload: ExchangeRatePayload{
				CountryName: "Nepal",
				Currency:    "Rupee",
			},
			want: ExchangeRateResponse{
				CountryCurrencyDesc: "Nepal-Rupee",
				ExchangeRate:        "130.5",
				RecordDate:          "2022-09-30",
			},
			wantErr: false,
		},
		{
			name:         "should return error for invalid country name",
			purchaseDate: "2022-11-30",
			payload: ExchangeRatePayload{
				CountryName: "Not",
				Currency:    "Rupee",
			},
			want:    ExchangeRateResponse{},
			wantErr: true,
		},
		{
			name:         "should return error for invalid country name",
			purchaseDate: "2022-11-30",
			payload: ExchangeRatePayload{
				CountryName: "Not",
				Currency:    "Rupee",
			},
			want:    ExchangeRateResponse{},
			wantErr: true,
		},
		{
			name:         "should return error for mismatch country and currency",
			purchaseDate: "2022-11-30",
			payload: ExchangeRatePayload{
				CountryName: "Nepal",
				Currency:    "Dollar",
			},
			want:    ExchangeRateResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExchangeRateGetter{}

			recordDate, err := time.Parse(time.DateOnly, tt.purchaseDate)
			if err != nil {
				t.Fatal()
			}

			got, err := e.GetExchangeRate(context.TODO(), ExchangeRatePayload{
				CountryName: tt.payload.CountryName,
				Currency:    tt.payload.Currency,
				RecordDate:  recordDate,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSixMonthBeforePurchaseDate(t *testing.T) {
	tests := []struct {
		given string
		want  string
		name  string
	}{
		{
			name:  "should successfully decrement year",
			given: "2023-06-30",
			want:  "2022-12-30",
		},
		{
			name:  "should correctly return 6 months before",
			given: "2023-02-28",
			want:  "2022-08-28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			givenDate, err := time.Parse(time.DateOnly, tt.given)
			if err != nil {
				t.Fatal()
			}
			sixMonthBefore := getSixMonthBeforePurchaseDate(givenDate)

			got := sixMonthBefore.Format(time.DateOnly)

			if got != tt.want {
				t.Errorf("getSixMonthBeforePurchaseDate() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestConvertAmount(t *testing.T) {
	type args struct {
		original     decimal.Decimal
		exchangeRate decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "Valid",
			args: args{
				original:     decimal.NewFromInt(100),
				exchangeRate: decimal.NewFromFloat(130.00),
			},
			want: decimal.NewFromFloat(13000.00),
		},
		{
			name: "Valid",
			args: args{
				original:     decimal.NewFromFloat(100),
				exchangeRate: decimal.NewFromFloat(132.90),
			},
			want: decimal.NewFromFloat(13290.00),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertAmount(tt.args.original, tt.args.exchangeRate)

			if !tt.want.Equal(got) {
				t.Errorf("Failed got = %s, want %s", got, tt.want)
			}
		})
	}
}
