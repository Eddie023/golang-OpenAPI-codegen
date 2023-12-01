package api

import (
	"log/slog"
	"net/http"

	"github.com/eddie023/wex-tag/ent/transaction"
	"github.com/eddie023/wex-tag/pkg/api/service"
	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/pkg/errors"
)

// GET api/v1/transaction/{transaction_id}
func (a *API) GetPurchaseTransaction(w http.ResponseWriter, r *http.Request, transactionId string, params types.GetPurchaseTransactionParams) {
	ctx := r.Context()

	uuidString, err := service.ParseStringToUUID(transactionId)
	if err != nil {
		slog.Error("failed to parse provided transaction id", "err", err.Error())
		apiout.Error(ctx, w, apiout.NewRequestError(errors.New("invalid transaction id provided"), http.StatusBadRequest))
		return
	}

	transaction, err := a.Db.Transaction.Query().Where(transaction.ID(uuidString)).First(ctx)
	if err != nil {
		apiout.Error(ctx, w, apiout.NewRequestError(errors.New("given transaction id not found"), http.StatusNotFound))
		return
	}

	exchangeInfo, err := a.ExchangeRateService.GetExchangeRate(ctx, service.ExchangeRatePayload{
		CountryName: params.Country,
		Currency:    params.Currency,
		RecordDate:  transaction.Date,
	})
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	output, err := a.ExchangeRateService.ConvertCurrency(service.ExchangeRatePayload{CountryName: params.Country, Currency: params.Currency}, transaction, exchangeInfo)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	apiout.JSON(ctx, w, output, http.StatusOK)
}

// POST api/v1/transaction
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
		apiout.Error(ctx, w, &apiout.BadRequestErr{
			Msg: err.Error(),
		})

		return
	}

	apiout.JSON(ctx, w, response, http.StatusCreated)
}
