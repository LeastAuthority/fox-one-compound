// Code generated by "stringer -type ActionType -trimprefix ActionType"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActionTypeDefault-0]
	_ = x[ActionTypeSupply-1]
	_ = x[ActionTypeBorrow-2]
	_ = x[ActionTypeRedeem-3]
	_ = x[ActionTypeRepay-4]
	_ = x[ActionTypeMint-5]
	_ = x[ActionTypePledge-6]
	_ = x[ActionTypeUnpledge-7]
	_ = x[ActionTypeLiquidate-8]
	_ = x[ActionTypeRedeemTransfer-9]
	_ = x[ActionTypeUnpledgeTransfer-10]
	_ = x[ActionTypeBorrowTransfer-11]
	_ = x[ActionTypeLiquidateTransfer-12]
	_ = x[ActionTypeRefundTransfer-13]
	_ = x[ActionTypeRepayRefundTransfer-14]
	_ = x[ActionTypeLiquidateRefundTransfer-15]
	_ = x[ActionTypeProposalAddMarket-16]
	_ = x[ActionTypeProposalUpdateMarket-17]
	_ = x[ActionTypeProposalWithdrawReserves-18]
	_ = x[ActionTypeProposalProvidePrice-19]
	_ = x[ActionTypeProposalVote-20]
	_ = x[ActionTypeProposalInjectCTokenForMint-21]
	_ = x[ActionTypeProposalUpdateMarketAdvance-22]
	_ = x[ActionTypeProposalTransfer-23]
	_ = x[ActionTypeProposalCloseMarket-24]
	_ = x[ActionTypeProposalOpenMarket-25]
	_ = x[ActionTypeProposalAddScope-26]
	_ = x[ActionTypeProposalRemoveScope-27]
	_ = x[ActionTypeProposalAddAllowList-28]
	_ = x[ActionTypeProposalRemoveAllowList-29]
	_ = x[ActionTypeUpdateMarket-30]
	_ = x[ActionTypeQuickPledge-31]
	_ = x[ActionTypeQuickBorrow-32]
	_ = x[ActionTypeQuickBorrowTransfer-33]
	_ = x[ActionTypeQuickRedeem-34]
	_ = x[ActionTypeQuickRedeemTransfer-35]
	_ = x[ActionTypeProposalAddOracleSigner-36]
	_ = x[ActionTypeProposalRemoveOracleSigner-37]
}

const _ActionType_name = "DefaultSupplyBorrowRedeemRepayMintPledgeUnpledgeLiquidateRedeemTransferUnpledgeTransferBorrowTransferLiquidateTransferRefundTransferRepayRefundTransferLiquidateRefundTransferProposalAddMarketProposalUpdateMarketProposalWithdrawReservesProposalProvidePriceProposalVoteProposalInjectCTokenForMintProposalUpdateMarketAdvanceProposalTransferProposalCloseMarketProposalOpenMarketProposalAddScopeProposalRemoveScopeProposalAddAllowListProposalRemoveAllowListUpdateMarketQuickPledgeQuickBorrowQuickBorrowTransferQuickRedeemQuickRedeemTransferProposalAddOracleSignerProposalRemoveOracleSigner"

var _ActionType_index = [...]uint16{0, 7, 13, 19, 25, 30, 34, 40, 48, 57, 71, 87, 101, 118, 132, 151, 174, 191, 211, 235, 255, 267, 294, 321, 337, 356, 374, 390, 409, 429, 452, 464, 475, 486, 505, 516, 535, 558, 584}

func (i ActionType) String() string {
	if i < 0 || i >= ActionType(len(_ActionType_index)-1) {
		return "ActionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActionType_name[_ActionType_index[i]:_ActionType_index[i+1]]
}
