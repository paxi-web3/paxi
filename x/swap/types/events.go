package types

const (
	EventParamsUpdated      = "ParamsUpdated"
	EventLiquidityProvided  = "LiquidityProvided"
	EventLiquidityWithdrawn = "LiquidityWithdrawn"
	EventSwapExecuted       = "SwapExecuted"

	AttributeKeyCreator      = "creator"
	AttributeKeyPrc20        = "prc20"
	AttributeKeyPaxiAmount   = "paxi_amount"
	AttributeKeyPrc20Amount  = "prc20_amount"
	AttributeKeyLPAmount     = "lp_amount"
	AttributeKeyLPMinted     = "lp_minted"
	AttributeKeyOfferDenom   = "offer_denom"
	AttributeKeyOfferAmount  = "offer_amount"
	AttributeKeyReceive      = "receive_amount"
	AttributeKeyMinReceive   = "min_receive"
	AttributeKeySwapSide     = "swap_side"
	AttributeKeyCodeID       = "code_id"
	AttributeKeySwapFeeBps   = "swap_fee_bps"
	AttributeKeyMinLiquidity = "min_liquidity"
	AttributeKeyReservePaxi  = "reserve_paxi"
	AttributeKeyReservePrc20 = "reserve_prc20"
	AttributeKeyTotalShares  = "total_shares"
	AttributeKeyPoolDeleted  = "pool_deleted"
)
