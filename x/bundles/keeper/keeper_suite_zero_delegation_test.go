package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - zero delegation

* Staker votes with zero delegation
* Staker receives vote slash with zero delegation
* Staker submit bundle proposal with zero delegation
* Staker receives upload slash with zero delegation
* Staker receives timeout slash because votes were missed
* Stakers try to produce valid bundle but all stakers have zero delegation

*/

var _ = Describe("valid bundles", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
	initialBalanceValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0)

	initialBalanceStaker1 := s.GetBalanceFromAddress(i.STAKER_1)
	initialBalanceValaddress1 := s.GetBalanceFromAddress(i.VALADDRESS_1)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  10_000,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Staker votes with zero delegation", func() {
		// ARRANGE
		// create normal validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// create zero delegation validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
	})

	It("Staker receives vote slash with zero delegation", func() {
		// ARRANGE
		// create normal validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// create zero delegation validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))

		// calculate voter slashes
		fraction, _ := sdk.NewDecFromStr(s.App().DelegationKeeper.GetVoteSlash(s.Ctx()))
		slashAmountVoter := uint64(sdk.NewDec(int64(0 * i.KYVE)).Mul(fraction).TruncateInt64())
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(0*i.KYVE - slashAmountVoter))

		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(100*i.KYVE - slashAmountVoter))
	})

	It("Staker submit bundle proposal with zero delegation", func() {
		// ARRANGE
		// create zero delegation validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		// create normal validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		valaccountUploader, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetBalanceFromAddress(valaccountUploader.Valaddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(valaccountUploader.Staker)

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// calculate uploader rewards
		totalReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64()) + pool.OperatingCost
		networkFee, _ := sdk.NewDecFromStr(s.App().BundlesKeeper.GetNetworkFee(s.Ctx()))

		treasuryReward := uint64(sdk.NewDec(int64(totalReward)).Mul(networkFee).TruncateInt64())
		totalUploaderReward := totalReward - treasuryReward

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0 + totalUploaderReward))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeZero())

		// check pool funds
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100*i.KYVE - totalReward))
	})

	It("Staker receives upload slash with zero delegation", func() {
		// ARRANGE
		// create zero delegation validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		// create normal validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.CommitAfterSeconds(60)

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentSummary).To(BeEmpty())
		Expect(pool.CurrentIndex).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		_, valaccountUploaderFound := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccountUploaderFound).To(BeFalse())

		balanceValaddress := s.GetBalanceFromAddress(i.VALADDRESS_0)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)
		_, uploaderFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(uploaderFound).To(BeTrue())

		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeZero())

		// calculate uploader slashes
		fraction, _ := sdk.NewDecFromStr(s.App().DelegationKeeper.GetUploadSlash(s.Ctx()))
		slashAmount := uint64(sdk.NewDec(int64(0 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(0*i.KYVE - slashAmount))
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(100*i.KYVE - slashAmount))

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeZero())

		// check pool funds
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100 * i.KYVE))
	})

	It("Staker receives timeout slash because votes were missed", func() {
		// ARRANGE
		// create normal validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// create zero delegation validator
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx()))

		for r := 1; r <= maxPoints; r++ {
			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.VALADDRESS_0,
				Staker:        i.STAKER_0,
				PoolId:        0,
				StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				DataSize:      100,
				DataHash:      "test_hash",
				FromIndex:     uint64(r * 100),
				BundleSize:    100,
				FromKey:       "test_key",
				ToKey:         "test_key",
				BundleSummary: "test_value",
			})

			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ASSERT
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, stakerFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_1)
		Expect(stakerFound).To(BeTrue())

		_, valaccountFound := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountFound).To(BeFalse())

		// check if voter got slashed
		slashAmountRatio, _ := sdk.NewDecFromStr(s.App().DelegationKeeper.GetTimeoutSlash(s.Ctx()))
		expectedBalance := 0*i.KYVE - uint64(sdk.NewDec(int64(0*i.KYVE)).Mul(slashAmountRatio).TruncateInt64())

		Expect(expectedBalance).To(Equal(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))
	})

	It("Stakers try to produce valid bundle but all stakers have zero delegation", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  0 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1)

		// manually set next staker
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		bundleProposal.NextUploader = i.STAKER_1
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentSummary).To(BeEmpty())
		Expect(pool.CurrentIndex).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		valaccountUploader, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccountUploader.Points).To(BeZero())

		balanceValaddress := s.GetBalanceFromAddress(valaccountUploader.Valaddress)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(valaccountUploader.Staker)

		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeZero())

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeZero())

		// check pool funds
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100 * i.KYVE))
	})
})
