package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRepayInsurance{}

type MsgRepayInsurance struct {
	Oracle sdk.AccAddress `json:"oracle"`
}

func NewMsgRepayInsurance(oracle sdk.AccAddress) MsgRepayInsurance {
	return MsgRepayInsurance{
		Oracle: oracle,
	}
}

// Indemnization relativo al valore da rimborsare
func (i Insurance) Indemnization(s Season) int64 {
	switch s {
	case Winter:
		return i.Money.Amount.Int64() * 12.0 / 10.0 // 120%
	case Spring:
		return i.Money.Amount.Int64() * 18.0 / 10.0 // 180%
	case Summer:
		return i.Money.Amount.Int64() * 30.0 / 10.0 // 300%
	default: // FALL
		return i.Money.Amount.Int64() * 15.0 / 10.0 // 150%
	}
}

const RepayInsuranceConst = "RepayInsurance"

func (msg MsgRepayInsurance) Route() string {
	return RouterKey
}

func (msg MsgRepayInsurance) Type() string {
	return RepayInsuranceConst
}

func (msg MsgRepayInsurance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Oracle}
}

func (msg MsgRepayInsurance) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRepayInsurance) ValidateBasic() error {
	if msg.Oracle.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Oracle.String())
	}
	return nil
}
