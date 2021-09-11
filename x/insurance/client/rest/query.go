package rest

import (
	"fmt"
	"net/http"

	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	route := fmt.Sprintf("/%s/%s", types.ModuleName, types.QueryAllInsurance)
	r.HandleFunc(route, queryInsurancesFn(cliCtx)).Methods("GET")
}

func queryInsurancesFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllInsurance)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
