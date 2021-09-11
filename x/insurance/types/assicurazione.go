package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
	"time"
)

const MIN = 1000
const MAX = 1000000000

type Insurance struct {
	ID           string         `json:"ID"`
	Intestatario sdk.AccAddress `json:"intestatario"`
	Money        sdk.Coin       `json:"money"`
	Day          int            `json:"day"`
	Month        int            `json:"month"`
	Year         int            `json:"year"`
	Rimborso     sdk.Coin       `json:"rimborso"`
}

// MonthInt converte il nome del mese in intero
func MonthInt(m time.Month) int {
	switch m.String() {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	default:
		return 12
	}
}

func (i Insurance) Validate() error {

	if i.ID == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "ID can't be empty")
	}

	if i.Intestatario.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, i.Intestatario.String())
	}

	return nil

}

func (i Insurance) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s,
												 intestatario: %s,
												 importo: %s,
												 giorno/mese/anno: %d/%d/%d,
												 rimborso: %s`,
		i.ID,
		i.Intestatario,
		i.Money,
		i.Day, i.Month, i.Year,
		i.Rimborso,
	))
}
