package keeper

import (
	"github.com/alessia19/insurance/x/insurance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "insurance_min_1000", InsuranceMin1000(k))
	ir.RegisterRoute(types.ModuleName, "insurance_max_1000000000", InsuranceMax1000000000(k))
}

// InsuranceMin1000 insurance deve avere un importo minimo di 1000
func InsuranceMin1000(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		insurances := keeper.GetAllInsurance(ctx)
		prezzoMinimo := sdk.NewCoin("foo", sdk.NewInt(1000))
		check := false
		for _, insurance := range insurances {
			check = check || insurance.Money.IsLT(prezzoMinimo)
		}

		return sdk.FormatInvariant(types.ModuleName,
			"insurance less than 1000 coins", "An insurance has an amount less than 1000 coins"), check
	}
}

func InsuranceMax1000000000(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		insurances := keeper.GetAllInsurance(ctx)
		prezzoMassimo := sdk.NewCoin("foo", sdk.NewInt(1000000000+1))
		check := false
		for _, insurance := range insurances {
			check = check || insurance.Money.IsGTE(prezzoMassimo)
		}

		return sdk.FormatInvariant(types.ModuleName,
				"insurance greater than 1000000000 coins",
				"An insurance has an amount greater than 1000000000 coins"),
			check
	}
}
