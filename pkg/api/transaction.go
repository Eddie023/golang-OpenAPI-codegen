package api

import (
	"net/http"

	"github.com/eddie023/wex-tag/pkg/apiout"
	"github.com/eddie023/wex-tag/pkg/types"
)

func (a *API) GetPurchaseTransaction(w http.ResponseWriter, r *http.Request, params types.GetPurchaseTransactionParams) {

}

func (a *API) PostPurchaseTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload types.CreateNewPurchaseTransaction

	err := apiout.DecodeJSONBody(w, r, &payload)
	if err != nil {
		apiout.Error(ctx, w, err)
		return
	}

	out, err := a.TransactionService.CreatePurchase(ctx, payload)
	if err != nil {
		apiout.Error(ctx, w, &apiout.BadRequestErr{
			Msg: "err",
		})
	}

	apiout.JSON(ctx, w, out, http.StatusCreated)
}
