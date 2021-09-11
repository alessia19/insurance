package keeper

import (
	"github.com/alessia19/insurance/x/insurance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_queryGetAllDebts(t *testing.T) {

	client, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")

	tests := []struct {
		name       string
		insurances []types.Insurance
	}{
		{
			"not insurance in storage",
			nil,
		},
		{
			"one insurance in storage",
			[]types.Insurance{
				{
					ID:           "1",
					Intestatario: client,
					Money:        sdk.NewCoin("foo", sdk.NewInt(78000)),
					Day:          23,
					Month:        4,
					Year:         2021,
					Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
				},
			},
		},
		{
			"more insurance in storage",
			[]types.Insurance{
				{
					ID:           "1",
					Intestatario: client,
					Money:        sdk.NewCoin("foo", sdk.NewInt(2000)),
					Day:          13,
					Month:        1,
					Year:         2021,
					Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
				},
				{
					ID:           "2",
					Intestatario: client,
					Money:        sdk.NewCoin("foo", sdk.NewInt(8007)),
					Day:          30,
					Month:        10,
					Year:         2021,
					Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
				},
				{
					ID:           "3",
					Intestatario: client,
					Money:        sdk.NewCoin("foo", sdk.NewInt(93400)),
					Day:          17,
					Month:        8,
					Year:         2021,
					Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, _, _, keeper := SetupTestInput()

			for _, insurance := range tt.insurances {
				require.NoError(t, keeper.BuyInsurance(ctx, insurance))
			}

			result, err := queryGetAllInsurance(ctx, nil, keeper)
			require.NoError(t, err)

			var ins []types.Insurance

			require.NotPanics(t, func() {
				cdc.MustUnmarshalJSON(result, &ins)
			})
			require.Equal(t, tt.insurances, ins)
		})
	}
}
