package claim

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/arkeonetwork/arkeo/testutil/sample"
	claimsimulation "github.com/arkeonetwork/arkeo/x/claim/simulation"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = claimsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

// nolint
const (
	opWeightMsgClaimEth = "op_weight_msg_claim_eth"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimEth int = 100

	opWeightMsgClaimArkeo = "op_weight_msg_claim_arkeo"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimArkeo int = 100

	opWeightMsgTransferClaim = "op_weight_msg_transfer_claim"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTransferClaim int = 100

	opWeightMsgAddClaim = "op_weight_msg_add_claim"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddClaim int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	claimGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&claimGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgClaimEth int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimEth, &weightMsgClaimEth, nil,
		func(_ *rand.Rand) {
			weightMsgClaimEth = defaultWeightMsgClaimEth
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimEth,
		claimsimulation.SimulateMsgClaimEth(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimArkeo int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimArkeo, &weightMsgClaimArkeo, nil,
		func(_ *rand.Rand) {
			weightMsgClaimArkeo = defaultWeightMsgClaimArkeo
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimArkeo,
		claimsimulation.SimulateMsgClaimArkeo(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTransferClaim int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTransferClaim, &weightMsgTransferClaim, nil,
		func(_ *rand.Rand) {
			weightMsgTransferClaim = defaultWeightMsgTransferClaim
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferClaim,
		claimsimulation.SimulateMsgTransferClaim(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddClaim int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddClaim, &weightMsgAddClaim, nil,
		func(_ *rand.Rand) {
			weightMsgAddClaim = defaultWeightMsgAddClaim
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddClaim,
		claimsimulation.SimulateMsgAddClaim(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
