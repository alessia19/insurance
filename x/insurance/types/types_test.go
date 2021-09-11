package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes_BuyDuration(t *testing.T) {
	client, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")

	amount := sdk.NewCoin("foo", sdk.NewInt(64000))
	day := 13
	month := 8
	year := 2021

	err1, buy1 := Buy(client, amount, day, month, year, 5)
	err2, buy2 := Buy(client, amount, day, month, year, 1)
	err3, buy3 := Buy(client, amount, day, month, year, 3)
	err4, buy4 := Buy(client, amount, day, month, year, 8)
	err5, buy5 := Buy(client, amount, day, month, year, 7)
	err6, buy6 := Buy(client, amount, day, month, year, 0)
	err7, buy7 := Buy(client, amount, day, month, year, 13)

	tests := []struct {
		name    string
		msgs    []sdk.Msg
		wantErr bool
	}{
		{
			"insurance for 5 days",
			buy1,
			false,
		},
		{
			"insurance for 1 day",
			buy2,
			false,
		},
		{
			"insurance for 3 days",
			buy3,
			false,
		},
		{
			"insurance for 8 days",
			buy4,
			true,
		},
		{
			"insurance for 7 days",
			buy5,
			false,
		},
		{
			"insurance for 0 days",
			buy6,
			true,
		},
		{
			"insurance for 13 days",
			buy7,
			true,
		},
	}

	var errors []error
	errors = append(errors, err1, err2, err3, err4, err5, err6, err7)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				require.Error(t, errors[i])
			} else {
				require.NoError(t, errors[i])
			}
		})
	}
}

func TestTypes_BuyAmount(t *testing.T) {
	client, _ := sdk.AccAddressFromBech32("cosmos1z56vw2lec6tkcyl5t8k2juqxtx74tg4pt88wmp")

	day := 13
	month := 8
	year := 2021

	err1, buy1 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(1000)), day, month, year, 5)
	err2, buy2 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(1000)), day, month, year, 1)
	err3, buy3 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(10000)), day, month, year, 3)
	err4, buy4 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(1950)), day, month, year, 2)
	err5, buy5 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(90000000000)), day, month, year, 7)
	err6, buy6 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(14000)), day, month, year, 4)
	err7, buy7 := Buy(client, sdk.NewCoin("foo", sdk.NewInt(2999)), day, month, year, 3)

	tests := []struct {
		name    string
		msgs    []sdk.Msg
		wantErr bool
	}{
		{
			"insurance for 5 days and 1000 coins of amount",
			buy1,
			true,
		},
		{
			"insurance for 1 day and 1000 coins of amount",
			buy2,
			false,
		},
		{
			"insurance for 3 days and 10000 coins of amount",
			buy3,
			false,
		},
		{
			"insurance for 2 days and 1950 coins of amount",
			buy4,
			true,
		},
		{
			"insurance for 7 days and 90000000000 coins of amount",
			buy5,
			true,
		},
		{
			"insurance for 4 days and 14000 coins of amount",
			buy6,
			false,
		},
		{
			"insurance for 3 days and 1950 coins of amount",
			buy7,
			true,
		},
	}

	var errors []error
	errors = append(errors, err1, err2, err3, err4, err5, err6, err7)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				require.Error(t, errors[i])
			} else {
				require.NoError(t, errors[i])
			}
		})
	}
}
