package cli

import (
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getAllInsurances(cdc),
	)

	return cmd
}

func getAllInsurances(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-insurances",
		Short: "Get all the insurances",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAllInsurancesFunc(cmd, args, cdc)
		},
	}
}

func getAllInsurancesFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllInsurance)
	res, _, err := cliCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}
