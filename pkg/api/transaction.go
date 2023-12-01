package api

import (
	"log/slog"
	"net/http"

	"github.com/eddie023/wex-tag/pkg/api/service"
	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/pkg/errors"
)

// GET /purchase/{transaction_id}?country=""&currency=""
func (a *API) GetPurchaseTransaction(w http.ResponseWriter, r *http.Request, transactionId string, params types.GetPurchaseTransactionParams) {
	ctx := r.Context()

	uuidString, err := service.ParseStringToUUID(transactionId)
	if err != nil {
		slog.Error("failed to parse provided transaction id", "err", err.Error())
		apiout.Error(ctx, w, apiout.NewRequestError(errors.New("invalid transaction id provided"), http.StatusBadRequest))
		return
	}

	transactionDetails, err := a.TransactionService.GetPurchaseDetailsByTransactionId(ctx, uuidString)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	exchangeRateDetails, err := a.ExchangeRateService.GetExchangeRate(ctx, service.ExchangeRatePayload{
		CountryName: params.Country,
		Currency:    params.Currency,
		RecordDate:  transactionDetails.Date,
	})
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	response, err := a.ExchangeRateService.ConvertCurrency(service.ExchangeRatePayload{CountryName: params.Country, Currency: params.Currency}, transactionDetails, exchangeRateDetails)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	apiout.JSON(ctx, w, response, http.StatusOK)
}

// POST /purchase
func (a *API) PostPurchaseTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload types.CreateNewPurchaseTransaction

	err := apiout.DecodeJSONBody(w, r, &payload)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	response, err := a.TransactionService.CreateNewPurchaseTransaction(ctx, payload)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	apiout.JSON(ctx, w, response, http.StatusCreated)
}
