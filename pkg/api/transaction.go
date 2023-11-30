package api

import (
	"net/http"

	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
)

// GET api/v1/transaction
func (a *API) GetPurchaseTransaction(w http.ResponseWriter, r *http.Request, params types.GetPurchaseTransactionParams) {

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

	response, err := a.TransactionService.CreateNewPurchase(ctx, payload)
	if err != nil {
		apiout.Error(ctx, w, &apiout.BadRequestErr{
			Msg: err.Error(),
		})

		return
	}

	apiout.JSON(ctx, w, response, http.StatusCreated)
}
