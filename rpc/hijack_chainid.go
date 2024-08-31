package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackChainID hijacks "eth_sendTransaction" method and sets the "chainID"
// field.
type hijackChainID struct {
	chainID uint64
	replace bool
}

func (c *hijackChainID) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		chainID := c.chainID
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if method != "eth_sendTransaction" || len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			// Verify arguments:
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return &ErrHijackFailed{name: "chain ID", err: fmt.Errorf("invalid transaction type: %T", args[0])}
			}

			// Get the transaction data:
			txd := tx.TransactionData()
			if !c.replace && txd.ChainID != nil {
				return next(ctx, t, result, method, args...)
			}

			// Get the chain ID from the RPC node:
			if chainID == 0 {
				chainID, err = (&MethodsCommon{Transport: t}).ChainID(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "chain ID", err: fmt.Errorf("failed to get chain ID: %w", err)}
				}
			}

			// Update the "chainID" field and continue:
			txd.ChainID = &chainID
			return next(ctx, t, result, method, args...)
		}
	}
}

func (c *hijackChainID) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (c *hijackChainID) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}
