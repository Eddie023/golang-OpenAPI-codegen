package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/api/mocks"
	"github.com/eddie023/wex-tag/pkg/api/service"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpMiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

/// NOTE:
//// These tests are primarily written to test the rigidity of our API layer rather than our business logic. Business logic test are written in service layer.

// newTestServer creates a configured API server for use in Go tests.
// The default time of the server is 1st Jan 2022, 10:00am UTC.
// This can be overriden by providing a custom clock with the withClock() option.
func newTestServer(t *testing.T, a *API) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(httpMiddleware.OapiRequestValidator(a.Swagger))
		types.HandlerWithOptions(a, types.ChiServerOptions{
			BaseRouter: r,
		})
	})

	return r
}

func TestPostTransactionAPI(t *testing.T) {
	type testcase struct {
		name                    string
		give                    string
		mockPurchaseTransaction *types.Transaction
		mockCreateErr           error

		wantCode int
		wantBody string
	}

	testDate, err := time.Parse(time.DateOnly, "2020-10-10")
	if err != nil {
		t.Fatal()
	}

	testUUID, err := uuid.Parse("680ed945-c2c3-4534-84e8-4ba6ed69eeea")
	if err != nil {
		t.Fatal()
	}

	testcases := []testcase{
		{
			name:     "should fail for emtpy json body with property missing error",
			give:     `{}`,
			wantCode: http.StatusBadRequest,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "123.45",
				Date:        time.Now().UTC(),
				Description: "abcd",
				Id:          uuid.New().String(),
			},
			mockCreateErr: nil,

			wantBody: `property "description" is missing`,
		},
		{
			name:     "should fail if amount field is not provided",
			give:     `{"description": ""}`,
			wantCode: http.StatusBadRequest,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "123.45",
				Date:        testDate,
				Description: "abcd",
				Id:          testUUID.String(),
			},
			mockCreateErr: nil,

			wantBody: `property "amount" is missing`,
		},
		{
			name:     "should successfully generate new purchase transaction details",
			give:     `{"description": "","amount": "1234.129123123123123123123213"}`,
			wantCode: http.StatusCreated,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "1234.129123123123123123123213",
				Date:        testDate.UTC(),
				Description: "",
				Id:          testUUID.String(),
			},
			mockCreateErr: nil,

			wantBody: `{"amountInUSD":"1234.129123123123123123123213","date":"2020-10-10T00:00:00Z","description":"","id":"680ed945-c2c3-4534-84e8-4ba6ed69eeea"}`,
		},
		{
			name:     "should fail with description cannot be longer than 50 chars",
			give:     `{"description": "text that is longer than 50 character text is longer than 50 character","amount": "1234.129123123123123123123213"}`,
			wantCode: http.StatusBadRequest,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "1234.129123123123123123123213",
				Date:        time.Now().UTC(),
				Description: "",
				Id:          "b6075c09-5fd3-4d6c-aa89-c9980d9be4d0",
			},
			mockCreateErr: nil,

			wantBody: `Error at "/description": maximum string length is 50`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockTransactionService(ctrl)
			if tc.mockPurchaseTransaction != nil {
				m.EXPECT().CreateNewPurchaseTransaction(gomock.Any(), gomock.Any()).Return(*tc.mockPurchaseTransaction, tc.mockCreateErr).AnyTimes()
			}

			swagger, err := types.GetSwagger()
			if err != nil {
				t.Fatal(err)
			}

			// remove any servers from the spec, as we don't know what host or port the user will run the API as.
			swagger.Servers = nil

			a := API{TransactionService: m, Swagger: swagger}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("POST", "/purchase", strings.NewReader(tc.give))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// this means we are trying to test fail condition
			if tc.wantCode != 200 && tc.wantCode != 201 {
				errorData, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Fatal(err)
				}

				if strings.Contains(string(errorData), tc.wantBody) {
					return
				}
			}

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(data) != tc.wantBody {
				t.Errorf("want =%s got =%s", tc.wantBody, string(data))
			}
		})
	}

}

func TestGetTransactionAPI(t *testing.T) {
	type testcase struct {
		name                string
		give                string
		queryParam          string
		mockExchangeRate    *service.ExchangeRateResponse
		mockExchangeRateErr error

		mockTransactionDetail *ent.Transaction

		wantCode int
		wantBody string
	}

	testDate, err := time.Parse(time.DateOnly, "2020-10-10")
	if err != nil {
		t.Fatal()
	}

	testUUID, err := uuid.Parse("680ed945-c2c3-4534-84e8-4ba6ed69eeea")
	if err != nil {
		t.Fatal()
	}

	testcases := []testcase{
		{
			name:             "should fail if country or currency query param is not passed",
			queryParam:       "?",
			give:             `{}`,
			wantCode:         http.StatusBadRequest,
			mockExchangeRate: &service.ExchangeRateResponse{},
			wantBody:         `parameter "country" in query has an error: value is required but missing`,
		},
		{
			name:             "should fail if only country param is passed",
			queryParam:       "country=Nepal",
			give:             `{}`,
			wantCode:         http.StatusBadRequest,
			mockExchangeRate: &service.ExchangeRateResponse{},
			wantBody:         `parameter "currency" in query has an error: value is required but missing`,
		},
		{
			name:             "should successfully return for valid query params",
			queryParam:       "country=Nepal&currency=Rupee",
			give:             `{}`,
			wantCode:         http.StatusOK,
			mockExchangeRate: &service.ExchangeRateResponse{},
			mockTransactionDetail: &ent.Transaction{
				ID:          testUUID,
				Date:        testDate,
				AmountInUsd: decimal.NewFromInt(100),
				Description: "",
			},
			wantBody: `{"convertedDetails":{"amount":"","country":"","currency":"","exchangeRateDate":"","exchangeRateUsed":""},"transactionDetails":{"amountInUSD":"","date":"0001-01-01T00:00:00Z","description":"","id":""}}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			exm := mocks.NewMockExchangeRateService(ctrl)
			if tc.mockExchangeRate != nil {
				exm.EXPECT().GetExchangeRate(gomock.Any(), gomock.Any()).Return(*tc.mockExchangeRate, tc.mockExchangeRateErr).AnyTimes()
				exm.EXPECT().ConvertCurrency(gomock.Any(), gomock.Any(), gomock.Any()).Return(types.GetPurchaseTransaction{}, nil).AnyTimes()
			}

			transm := mocks.NewMockTransactionService(ctrl)
			if tc.mockTransactionDetail != nil {
				transm.EXPECT().GetPurchaseDetailsByTransactionId(gomock.Any(), gomock.Any()).Return(tc.mockTransactionDetail, nil).AnyTimes()
			}

			swagger, err := types.GetSwagger()
			if err != nil {
				t.Fatal(err)
			}

			// remove any servers from the spec, as we don't know what host or port the user will run the API as.
			swagger.Servers = nil

			a := API{ExchangeRateService: exm, TransactionService: transm, Swagger: swagger}
			handler := newTestServer(t, &a)

			req, err := http.NewRequest("GET", fmt.Sprintf("/purchase/%s?%s", testUUID, tc.queryParam), strings.NewReader(tc.give))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// this means we are trying to test fail condition
			if tc.wantCode != 200 && tc.wantCode != 201 {
				errorData, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Fatal(err)
				}

				if strings.Contains(string(errorData), tc.wantBody) {
					return
				} else {
					t.Errorf("want =%s got=%s", tc.wantBody, errorData)
					return
				}
			}

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(data) != tc.wantBody {
				t.Errorf("want =%s got =%s", tc.wantBody, string(data))
			}

		})
	}

}
