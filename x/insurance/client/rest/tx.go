package rest

// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/buyinsurance", types.ModuleName), buyInsuranceFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/repayinsurance", types.ModuleName), repayFn(cliCtx)).Methods("POST")
}

// Action TX body
type buyInsuranceRequest struct {
	BaseReq      rest.BaseReq   `json:"base_req"` // i dati di una richiesta Cosmos di base
	ID           string         `json:"ID"`
	Intestatario sdk.AccAddress `json:"intestatario"`
	Money        sdk.Coin       `json:"money"`
	Day          int            `json:"day"`
	Month        int            `json:"month"`
	Year         int            `json:"year"`
	Rimborso     sdk.Coin       `json:"rimborso"`
}

type repayInsuranceRequest struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Oracle  sdk.AccAddress `json:"oracle"`
}

func buyInsuranceFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyInsuranceRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		// create the message
		msg := types.MsgBuyInsurance{
			ID:           req.ID,
			Intestatario: req.Intestatario,
			Money:        req.Money,
			Day:          req.Day,
			Month:        req.Month,
			Year:         req.Year,
			Rimborso:     req.Rimborso,
		}

		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func repayFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var repay repayInsuranceRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &repay) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := repay.BaseReq.Sanitize()

		if !baseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		// create the message
		msg := types.MsgRepayInsurance{
			Oracle: repay.Oracle,
		}

		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
