package insurance

import (
	"github.com/alessia19/insurance/x/insurance/keeper"
	"github.com/alessia19/insurance/x/insurance/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper       = keeper.NewKeeper
	NewQuerier      = keeper.NewQuerier
	RegisterCodec   = types.RegisterCodec
	NewGenesisState = types.NewGenesisState
	ValidateGenesis = types.ValidateGenesis

	// variabile aliases
	ModuleCdc = types.ModuleCdc

	NewMsgBuyInsurance   = types.NewMsgBuyInsurance
	NewMsgRepayInsurance = types.NewMsgRepayInsurance
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgBuyInsurance   = types.MsgBuyInsurance
	MsgRepayInsurance = types.MsgRepayInsurance
)
