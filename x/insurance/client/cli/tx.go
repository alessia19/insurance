package cli

import (
	"bufio"
	"fmt"
	"github.com/alessia19/insurance/x/insurance/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		buyInsuranceCmd(cdc),
		repayInsuranceCmd(cdc),
	)
	return txCmd
}

func buyInsuranceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [dd/mm/yyyy] [duration] [amount]",
		Short: "Buy an insurance",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return buyInsuranceCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func buyInsuranceCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	intestatario := cliCtx.GetFromAddress()

	date := args[0]
	d := strings.Split(date, "/")

	day, err := strconv.Atoi(d[0])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(d[1])
	if err != nil {
		return err
	}
	year, err := strconv.Atoi(d[2])
	if err != nil {
		return err
	}

	duration, _ := strconv.Atoi(args[1])

	amount, _ := sdk.ParseCoin(args[2])

	nameIntestatario := cliCtx.GetFromName()
	// controlla che il nome non corrispondi a "oracle"
	if strings.Compare(nameIntestatario, "oracle") == 0 {
		return fmt.Errorf("the name of the intestatario cannot be %s", nameIntestatario)
	}

	err, msgs := types.Buy(intestatario, amount, day, month, year, duration)

	if err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
}

func repayInsuranceCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repay",
		Short: "Refund insured insurance for today's date",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return repayInsuranceCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func repayInsuranceCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	oracle := cliCtx.GetFromAddress() // indirizzo del messaggio repay
	name := cliCtx.GetFromName()      // nome dell'indirizzo di questo messaggio

	if strings.Compare(name, "oracle") != 0 {
		return fmt.Errorf("the name isn't \"oracle\", but \"%s\"", name)
	}

	msg := types.NewMsgRepayInsurance(oracle)

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
