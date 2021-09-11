package keeper

import (
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for insurance clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryAllInsurance:
			return queryGetAllInsurance(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetAllInsurance(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, keeper.GetAllInsurance(ctx))
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
