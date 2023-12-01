package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eddie023/wex-tag/pkg/api/mocks"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpMiddleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

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

	testcases := []testcase{
		{
			name:     "emtpy body should fail with property missing",
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
			name:     "amount field is required for request body",
			give:     `{"description": ""}`,
			wantCode: http.StatusBadRequest,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "123.45",
				Date:        time.Now().UTC(),
				Description: "abcd",
				Id:          uuid.New().String(),
			},
			mockCreateErr: nil,

			wantBody: `property "amount" is missing`,
		},
		{
			name:     "ok 2",
			give:     `{"description": "","amount": "1234.129123123123123123123213"}`,
			wantCode: http.StatusCreated,

			mockPurchaseTransaction: &types.Transaction{
				AmountInUSD: "1234.129123123123123123123213",
				Date:        time.Now().UTC(),
				Description: "",
				Id:          "b6075c09-5fd3-4d6c-aa89-c9980d9be4d0",
			},
			mockCreateErr: nil,

			wantBody: `{"amountInUSD":"1234.129123123123123123123213","date":"2023-11-29 09:20:03.043085 +0000 UTC","description":"","id":"b6075c09-5fd3-4d6c-aa89-c9980d9be4d0"}`,
		},
		{
			name:     "description cannot be longer than 50 chars",
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

			req, err := http.NewRequest("POST", "/transaction", strings.NewReader(tc.give))
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
		name string
		give string

		wantCode int
		wantBody string
	}

	testcases := []testcase{
		{
			name:     "emtpy body should fail with property missing",
			give:     `{}`,
			wantCode: http.StatusBadRequest,
			wantBody: `property "description" is missing`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

		})
	}

}
