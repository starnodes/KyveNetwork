package keeper

import (
	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetVestingStatus returns all computed values which are dependent on the time
// for the given account
func GetVestingStatus(account types.TeamVestingAccount, time uint64) *types.VestingStatus {
	status := types.VestingStatus{}

	// get total allocation
	status.TotalVestedAmount = getVestedAmount(account, time)
	status.TotalUnlockedAmount = getUnlockedAmount(account, time)
	if status.TotalUnlockedAmount > account.UnlockedClaimed {
		status.CurrentClaimableAmount = getUnlockedAmount(account, time) - account.UnlockedClaimed
	}

	status.LockedVestedAmount = status.TotalVestedAmount - status.TotalUnlockedAmount
	status.RemainingUnvestedAmount = getVestingMaxAmount(account) - status.TotalVestedAmount

	return &status
}

// GetVestingPlan returns all computed static values for a given account
func GetVestingPlan(account types.TeamVestingAccount) *types.VestingPlan {
	plan := types.VestingPlan{}

	plan.MaximumVestingAmount = getVestingMaxAmount(account)
	plan.ClawbackAmount = account.TotalAllocation - plan.MaximumVestingAmount

	plan.TokenVestingStart = account.Commencement + types.CLIFF_DURATION
	plan.TokenVestingFinished = account.Commencement + types.VESTING_DURATION

	plan.TokenUnlockStart = getLockUpReferenceDate(account)
	plan.TokenUnlockFinished = getLockUpReferenceDate(account) + types.UNLOCK_DURATION

	return &plan
}

// GetIssuedTeamAllocation gets the total amount in $KYVE which is issued to all team vesting accounts.
// It is equal to the sum of all max vesting amounts, because normally the usage of all
// vesting accounts is the sum of all allocations minus the clawback which getVestingMaxAmount
// already takes into account
func (k Keeper) GetIssuedTeamAllocation(ctx sdk.Context) (used uint64) {
	for _, account := range k.GetTeamVestingAccounts(ctx) {
		used += getVestingMaxAmount(account)
	}

	return
}

func (k Keeper) GetTeamInfo(ctx sdk.Context) (info *types.QueryTeamInfoResponse) {
	info = &types.QueryTeamInfoResponse{}

	info.Authority = types.AUTHORITY_ADDRESS
	info.TotalTeamAllocation = types.TEAM_ALLOCATION

	info.IssuedTeamAllocation = k.GetIssuedTeamAllocation(ctx)
	info.AvailableTeamAllocation = types.TEAM_ALLOCATION - info.IssuedTeamAllocation

	authority := k.GetAuthority(ctx)
	info.TotalAuthorityRewards = authority.TotalRewards
	info.ClaimedAuthorityRewards = authority.RewardsClaimed
	info.AvailableAuthorityRewards = authority.TotalRewards - authority.RewardsClaimed

	info.RequiredModuleBalance = types.TEAM_ALLOCATION + info.AvailableAuthorityRewards

	for _, account := range k.GetTeamVestingAccounts(ctx) {
		info.TotalAccountRewards += account.TotalRewards
		info.ClaimedAccountRewards += account.RewardsClaimed
		info.AvailableAccountRewards += account.TotalRewards - account.RewardsClaimed

		info.RequiredModuleBalance += account.TotalRewards - account.RewardsClaimed
		info.RequiredModuleBalance -= account.UnlockedClaimed
	}

	coins := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), globalTypes.Denom)
	info.TeamModuleBalance = uint64(coins.Amount.Int64())

	return
}

// getVestedAmount returns the total amount of $KYVE that has vested until the given time for the given user.
// The function is well-defined for all values of t
func getVestedAmount(account types.TeamVestingAccount, time uint64) uint64 {
	// the account vesting duration is the time in seconds an account is already vesting
	accountVestingDuration := uint64(0)

	// if the specified time is after the commencement date the vesting duration is longer than zero
	if time > account.Commencement {
		accountVestingDuration = time - account.Commencement
	}

	// if a clawback time is defined and if it is before the specified time the vesting duration only goes
	// until the clawback time
	if account.Clawback > 0 && account.Clawback < time {
		accountVestingDuration = account.Clawback - account.Commencement
	}

	// if account is vesting less than the vesting cliff the vested amount is zero
	if accountVestingDuration < types.CLIFF_DURATION {
		return 0
	}

	// if user is vesting less than the vesting duration the vested amount is linear to the membership time
	if accountVestingDuration < types.VESTING_DURATION {
		vested := sdk.NewDec(int64(account.TotalAllocation)).
			Mul(sdk.NewDec(int64(accountVestingDuration))).
			Quo(sdk.NewDec(int64(types.VESTING_DURATION)))

		return uint64(vested.TruncateInt64())
	}

	// if user is vesting longer than the vesting duration the entire allocation has vested
	return account.TotalAllocation
}

// getVestingMaxAmount gets the maximum amount an account can possibly vest
func getVestingMaxAmount(account types.TeamVestingAccount) uint64 {
	// in order to get the maximum possible vesting amount we add the total vesting duration to the
	// commencement date as the specified time
	return getVestedAmount(account, account.Commencement+types.VESTING_DURATION)
}

// getLockUpReferenceDate gets the unix time the unlocking starts for an account
func getLockUpReferenceDate(account types.TeamVestingAccount) uint64 {
	// the unlocking starts exactly 1 year after the commencement or TGE, whatever the latter is
	return util.MaxUInt64(account.Commencement, types.TGE) + types.CLIFF_DURATION
}

// getUnlockedAmount returns total amount of $KYVE that has unlocked until the given time for the given user.
// The function is well-defined for all values of t
func getUnlockedAmount(account types.TeamVestingAccount, time uint64) uint64 {
	// get the unix time the unlocking for an account starts
	timeUnlock := getLockUpReferenceDate(account)

	// if the specified time is before the lockup reference data the unlocked amount is zero
	if time < timeUnlock {
		return 0
	}
	// => time - timeUnlock >= 0

	if time-timeUnlock < types.UNLOCK_DURATION {
		// get the total vested amount based on specified time
		vested := getVestedAmount(account, time)

		// calculate the unlocked amount linearly based on time
		unlocked := sdk.NewDec(int64(vested)).
			Mul(sdk.NewDec(int64(time - timeUnlock))).
			Quo(sdk.NewDec(int64(types.UNLOCK_DURATION)))

		return uint64(unlocked.TruncateInt64())
	}

	// if specified time comes after the unlock duration the full maximum vesting amount is unlocked
	return getVestingMaxAmount(account)
}
