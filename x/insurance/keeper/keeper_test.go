package keeper

import (
	"github.com/alessia19/insurance/x/insurance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestKeeper_CreateInsurance crea assicurazioni per lo stesso cliente
func TestKeeper_CreateInsurance(t *testing.T) {
	intestatario, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")

	insurance1 := types.Insurance{
		ID:           "0001",
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(34123)),
		Day:          8,
		Month:        8,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	_, ctx, _, _, keeper := SetupTestInput()
	err := keeper.BuyInsurance(ctx, insurance1)
	require.NoError(t, err)

	insurance2 := types.Insurance{
		ID:           "0002",
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(1200)),
		Day:          8,
		Month:        6,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	err = keeper.BuyInsurance(ctx, insurance2)
	require.NoError(t, err)

	insurance3 := types.Insurance{
		ID:           "5",
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(356000)),
		Day:          5,
		Month:        10,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	err = keeper.BuyInsurance(ctx, insurance3)
	require.NoError(t, err)
}

// TestKeeper_CreateInsuranceWithSameID anche se due insurance vengono create con lo stesso ID
// l'insurance successiva avrà un ID distinto e non uguale al precedente
func TestKeeper_CreateInsuranceWithSameID(t *testing.T) {
	intestatario, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")

	insurance1 := types.Insurance{
		ID:           "0006",
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(5300)),
		Day:          9,
		Month:        5,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	_, ctx, _, _, keeper := SetupTestInput()
	err := keeper.BuyInsurance(ctx, insurance1)
	require.NoError(t, err)

	insurance2 := types.Insurance{
		ID:           "0006", // stesso identificatore
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(5300)),
		Day:          9,
		Month:        5,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}
	err = keeper.BuyInsurance(ctx, insurance2)
	require.NoError(t, err)

	insurance3 := types.Insurance{
		ID:           "0006",
		Intestatario: intestatario,
		Money:        sdk.NewCoin("foo", sdk.NewInt(5300)),
		Day:          3,
		Month:        8,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	err = keeper.BuyInsurance(ctx, insurance3)
	require.NoError(t, err)
}

// TestKeeper_CreateInsuranceForDifferentClients crea un'assicurazione per ogni cliente
func TestKeeper_CreateInsuranceForDifferentClients(t *testing.T) {
	client1, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")
	client2, _ := sdk.AccAddressFromBech32("cosmos18670tqdkvqeyfjx5g4txll4s0sfk73s67sycat")

	insurance1 := types.Insurance{
		ID:           "10",
		Intestatario: client1,
		Money:        sdk.NewCoin("foo", sdk.NewInt(3000)),
		Day:          15,
		Month:        7,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	_, ctx, _, _, keeper := SetupTestInput()
	err := keeper.BuyInsurance(ctx, insurance1)

	require.NoError(t, err)

	insurance2 := types.Insurance{
		ID:           "11",
		Intestatario: client2,
		Money:        sdk.NewCoin("foo", sdk.NewInt(13000)),
		Day:          24,
		Month:        4,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	err = keeper.BuyInsurance(ctx, insurance2)
	require.NoError(t, err)
}

func TestKeeper_RepayInsurance(t *testing.T) {
	client, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")
	oracle, _ := sdk.AccAddressFromBech32("cosmos18670tqdkvqeyfjx5g4txll4s0sfk73s67sycat")

	insurance1 := types.Insurance{
		ID:           "1",
		Intestatario: client,
		Money:        sdk.NewCoin("foo", sdk.NewInt(1200)),
		Day:          time.Now().Day(),
		Month:        types.MonthInt(time.Now().Month()),
		Year:         time.Now().Year(),
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	_, ctx, _, bankKeeper, keeper := SetupTestInput()
	err := keeper.BuyInsurance(ctx, insurance1)
	require.NoError(t, err)

	insurance2 := types.Insurance{
		ID:           "2",
		Intestatario: client,
		Money:        sdk.NewCoin("foo", sdk.NewInt(2300)),
		Day:          12,
		Month:        4,
		Year:         2021,
		Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
	}

	err = keeper.BuyInsurance(ctx, insurance2)
	require.NoError(t, err)

	repay := types.MsgRepayInsurance{
		Oracle: oracle,
	}

	var insurances []types.Insurance
	insurances = append(insurances, insurance1, insurance2)

	// imposta le monete dell'oracolo
	require.NoError(t, bankKeeper.SetCoins(ctx, repay.Oracle, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(10000)))))

	err = keeper.RepayInsurance(ctx, repay)
	require.NoError(t, err)

	var newRimborso sdk.Coin

	for _, ins := range insurances {
		if ins.Day == time.Now().Day() && ins.Month == types.MonthInt(time.Now().Month()) &&
			ins.Year == time.Now().Year() {
			newRimborso = sdk.NewCoin(ins.Rimborso.Denom, sdk.NewInt(ins.Indemnization(types.Now())))
			ins.Rimborso = ins.Rimborso.Add(newRimborso)
		} else {
			newRimborso = sdk.NewCoin(ins.Rimborso.Denom, sdk.NewInt(0))
		}

		require.True(t, ins.Rimborso.IsEqual(newRimborso))
	}
}

func TestKeeper_RepayInsurance2(t *testing.T) {
	client, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")
	oracle, _ := sdk.AccAddressFromBech32("cosmos18670tqdkvqeyfjx5g4txll4s0sfk73s67sycat")

	tests := []struct {
		name          string
		assicurazione *types.Insurance
		wantErr       bool
	}{
		{
			"insurance 1",
			&types.Insurance{
				ID:           "1",
				Intestatario: client,
				Money:        sdk.NewCoin("foo", sdk.NewInt(1000)),
				Day:          time.Now().Day(),
				Month:        types.MonthInt(time.Now().Month()),
				Year:         time.Now().Year(),
				Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
			},
			false,
		},
		{
			"insurance 2",
			&types.Insurance{
				ID:           "2",
				Intestatario: client,
				Money:        sdk.NewCoin("foo", sdk.NewInt(23000)),
				Day:          9,
				Month:        12,
				Year:         2021,
				Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
			},
			true,
		},
		{
			"insurance 3",
			&types.Insurance{
				ID:           "3",
				Intestatario: client,
				Money:        sdk.NewCoin("foo", sdk.NewInt(9123)),
				Day:          29,
				Month:        7,
				Year:         2021,
				Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
			},
			true,
		},
		{
			"insurance 4",
			&types.Insurance{
				ID:           "4",
				Intestatario: client,
				Money:        sdk.NewCoin("foo", sdk.NewInt(84001)),
				Day:          time.Now().Day(),
				Month:        types.MonthInt(time.Now().Month()),
				Year:         time.Now().Year(),
				Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
			},
			false,
		},
		{
			"insurance 5",
			&types.Insurance{
				ID:           "5",
				Intestatario: client,
				Money:        sdk.NewCoin("foo", sdk.NewInt(120345)),
				Day:          9,
				Month:        4,
				Year:         2021,
				Rimborso:     sdk.NewCoin("foo", sdk.NewInt(0)),
			},
			true,
		},
	}

	repay := types.MsgRepayInsurance{
		Oracle: oracle,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, _, bankKeeper, keeper := SetupTestInput()
			// se c'è un'assicurazione
			if tt.assicurazione != nil {
				require.NoError(t, tt.assicurazione.Validate())
				// creo l'assicurazione in memoria
				require.NoError(t, keeper.BuyInsurance(ctx, *tt.assicurazione))
			}

			require.NoError(t, bankKeeper.SetCoins(ctx, repay.Oracle, sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(10000)))))
			err := keeper.RepayInsurance(ctx, repay)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			// altrimenti se non ci sono errori
			newRimborso := sdk.NewCoin(tt.assicurazione.Rimborso.Denom, sdk.NewInt(tt.assicurazione.Indemnization(types.Now())))
			tt.assicurazione.Rimborso = tt.assicurazione.Rimborso.Add(newRimborso)

			require.True(t, tt.assicurazione.Rimborso.IsEqual(newRimborso))
		})
	}
}
