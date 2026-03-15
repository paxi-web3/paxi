package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidRequest      = errorsmod.Register(ModuleName, 1, "invalid request")
	ErrMarketNotFound      = errorsmod.Register(ModuleName, 2, "market not found")
	ErrPositionNotFound    = errorsmod.Register(ModuleName, 3, "position not found")
	ErrInvalidMarketStatus = errorsmod.Register(ModuleName, 4, "invalid market status")
	ErrInvalidOutcome      = errorsmod.Register(ModuleName, 5, "invalid outcome")
	ErrInsufficientFunds   = errorsmod.Register(ModuleName, 6, "insufficient funds")
	ErrUnauthorized        = errorsmod.Register(ModuleName, 7, "unauthorized")
	ErrDuplicateTrade      = errorsmod.Register(ModuleName, 8, "duplicate trade id")
	ErrAlreadyClaimed      = errorsmod.Register(ModuleName, 9, "already claimed")
	ErrOrderNotFound       = errorsmod.Register(ModuleName, 10, "order not found")
	ErrInvalidOrderStatus  = errorsmod.Register(ModuleName, 11, "invalid order status")
	ErrInvalidOrderPair    = errorsmod.Register(ModuleName, 12, "invalid order pair")
)
