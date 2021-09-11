package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"time"
)

// verify interface at compile time
var _ sdk.Msg = &MsgBuyInsurance{}

type MsgBuyInsurance Insurance

func NewMsgBuyInsurance(i Insurance) MsgBuyInsurance {
	return MsgBuyInsurance(i)
}

const BuyInsuranceConst = "BuyInsurance"

func (msg MsgBuyInsurance) Route() string { return RouterKey }
func (msg MsgBuyInsurance) Type() string  { return BuyInsuranceConst }
func (msg MsgBuyInsurance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Intestatario}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgBuyInsurance) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Buy assicura dalla data passata per duration giorni
func Buy(intestatario sdk.AccAddress, amount sdk.Coin, day, month, year, duration int) (error, []sdk.Msg) {

	ID := 1
	var n string // variabile per contenere la stringa dell'ID
	n = strconv.Itoa(ID)
	for i := 0; i < 5; i++ {
		if len(n) != 4 {
			n = "0" + n
		}
	}

	// data d'inizio dell'assicurazione
	data := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	if duration < 1 {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "you must insure at least one day"), nil
	}

	if duration > 7 {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "you cannot insure more than one week"), nil
	}

	if amount.Amount.LT(sdk.NewInt(MIN * int64(duration))) {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("we insure a single day for at least %d units of coin", MIN)), nil
	}

	if amount.Amount.GT(sdk.NewInt(MAX * int64(duration))) {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("we insure a single day for up to %d units of coin", MAX)), nil
	}

	var msgs []sdk.Msg

	for offset := 0; offset < duration; offset++ {
		d := data.AddDate(0, 0, offset)

		msg := MsgBuyInsurance(Insurance{
			ID:           n,
			Intestatario: intestatario,
			Money:        sdk.NewCoin(amount.Denom, sdk.NewInt(amount.Amount.Int64()/int64(duration))),
			Day:          d.Day(),
			Month:        MonthInt(d.Month()),
			Year:         d.Year(),
			Rimborso:     sdk.NewCoin(amount.Denom, sdk.NewInt(0)),
		})

		if err := msg.ValidateBasic(); err != nil {
			return err, nil
		}

		msgs = append(msgs, msg)
	}
	return nil, msgs
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgBuyInsurance) ValidateBasic() error {
	return Insurance(msg).Validate()
}
