package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Simulation parameter constants
const (
	UnbondingTime = "unbonding_time"
	MaxValidators = "max_validators"
)

// GenUnbondingTime randomized UnbondingTime
func GenUnbondingTime(cdc *codec.Codec, r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// GenMaxValidators randomized MaxValidators
func GenMaxValidators(cdc *codec.Codec, r *rand.Rand) (maxValidators uint16) {
	return uint16(r.Intn(250) + 1)
}

// GenStakingGenesisState generates a random GenesisState for staking
func GenStakingGenesisState(
	cdc *codec.Codec, r *rand.Rand, genesisState map[string]json.RawMessage,
	accs []simulation.Account, amount, numAccs, numInitiallyBonded int64,
) staking.GenesisState {

	var (
		validators  []staking.Validator
		delegations []staking.Delegation
	)

	// params
	ubdTime := GenUnbondingTime(cdc, r)
	maxValidators := GenMaxValidators(cdc, r)
	params := staking.NewParams(ubdTime, maxValidators, 7, sdk.DefaultBondDenom)

	// validators & delegations
	valAddrs := make([]sdk.ValAddress, numInitiallyBonded)
	for i := 0; i < int(numInitiallyBonded); i++ {
		valAddr := sdk.ValAddress(accs[i].Address)
		valAddrs[i] = valAddr

		validator := staking.NewValidator(valAddr, accs[i].PubKey, staking.Description{})
		validator.Tokens = sdk.NewInt(amount)
		validator.DelegatorShares = sdk.NewDec(amount)
		delegation := staking.NewDelegation(accs[i].Address, valAddr, sdk.NewDec(amount))
		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := staking.NewGenesisState(params, validators, delegations)

	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", codec.MustMarshalJSONIndent(cdc, stakingGenesis.Params))
	genesisState[staking.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)

	return stakingGenesis
}
