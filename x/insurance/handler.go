package insurance

import (
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the insurance type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgBuyInsurance:
			return handleMsgBuyInsurance(ctx, keeper, msg)
		case types.MsgRepayInsurance:
			return handleMsgRepayInsurance(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgBuyInsurance(ctx sdk.Context, k Keeper, msg MsgBuyInsurance) (*sdk.Result, error) {
	err := k.BuyInsurance(ctx, types.Insurance(msg))

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Insurance creata con successo"}, nil
}

func handleMsgRepayInsurance(ctx sdk.Context, k Keeper, msg MsgRepayInsurance) (*sdk.Result, error) {
	err := k.RepayInsurance(ctx, msg)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Insurance rimborsata con successo"}, nil
}
