package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackLegacyGasFee hijacks "eth_sendTransaction" method and sets the
// "gasPrice" field using the estimate provided by the RPC node.
type hijackLegacyGasFee struct {
	multiplier      float64
	minGasPrice     *big.Int
	maxGasPrice     *big.Int
	replace         bool
	allowChangeType bool
}

func (h *hijackLegacyGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if method != "eth_sendTransaction" || len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			// Verify arguments:
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return &ErrHijackFailed{name: "legacy gas price", err: fmt.Errorf("invalid transaction type: %T", args[0])}
			}

			// Change transaction type if needed:
			if h.allowChangeType {
				tx = convertTXToLegacyPrice(tx)
				args[0] = tx
			}

			// Get legacy price data:
			lpd := getLegacyPriceData(tx)
			if lpd == nil {
				return next(ctx, t, result, method, args...)
			}

			// Skip if already set:
			if !h.replace && lpd.GasPrice != nil {
				return next(ctx, t, result, method, args...)
			}

			// Get gas price from RPC node:
			gasPrice, err := (&MethodsCommon{Transport: t}).GasPrice(ctx)
			if err != nil {
				return &ErrHijackFailed{name: "legacy gas price", err: fmt.Errorf("failed to get gas price: %w", err)}
			}

			// Update gas price and continue:
			gasPrice, _ = new(big.Float).Mul(new(big.Float).SetInt(gasPrice), big.NewFloat(h.multiplier)).Int(nil)
			if h.minGasPrice != nil && gasPrice.Cmp(h.minGasPrice) < 0 {
				gasPrice = h.minGasPrice
			}
			if h.maxGasPrice != nil && gasPrice.Cmp(h.maxGasPrice) > 0 {
				gasPrice = h.maxGasPrice
			}
			lpd.GasPrice = gasPrice
			return next(ctx, t, result, method, args...)
		}
	}
}

func (h *hijackLegacyGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (h *hijackLegacyGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

// hijackDynamicGasFee hijacks "eth_sendTransaction" method and sets the
// "maxFeePerGas" and "maxPriorityFeePerGas" fields using the estimate provided
// by the RPC node.
type hijackDynamicGasFee struct {
	gasPriceMultiplier          float64
	priorityFeePerGasMultiplier float64
	minGasPrice                 *big.Int
	maxGasPrice                 *big.Int
	minPriorityFeePerGas        *big.Int
	maxPriorityFeePerGas        *big.Int
	replace                     bool
	allowChangeType             bool
}

func (h *hijackDynamicGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if method != "eth_sendTransaction" || len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			// Verify arguments:
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("invalid transaction type: %T", args[0])}
			}

			// Change transaction type if needed:
			if h.allowChangeType {
				tx = convertTXToDynamicFee(tx)
				args[0] = tx
			}

			// Get dynamic fee data:
			dfd := getDynamicFeeData(tx)
			if dfd == nil {
				return next(ctx, t, result, method, args...)
			}

			// Skip if already set:
			if !h.replace && dfd.MaxFeePerGas != nil && dfd.MaxPriorityFeePerGas != nil {
				return next(ctx, t, result, method, args...)
			}

			// Get gas price from RPC node:
			maxFeePerGas, err := (&MethodsCommon{Transport: t}).GasPrice(ctx)
			if err != nil {
				return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get gas price: %w", err)}
			}
			priorityFeePerGas, err := (&MethodsCommon{Transport: t}).MaxPriorityFeePerGas(ctx)
			if err != nil {
				return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get priority fee per gas: %w", err)}
			}

			// Update gas price and continue:
			maxFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(maxFeePerGas), big.NewFloat(h.gasPriceMultiplier)).Int(nil)
			priorityFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(priorityFeePerGas), big.NewFloat(h.priorityFeePerGasMultiplier)).Int(nil)
			if h.minGasPrice != nil && maxFeePerGas.Cmp(h.minGasPrice) < 0 {
				maxFeePerGas = h.minGasPrice
			}
			if h.maxGasPrice != nil && maxFeePerGas.Cmp(h.maxGasPrice) > 0 {
				maxFeePerGas = h.maxGasPrice
			}
			if h.minPriorityFeePerGas != nil && priorityFeePerGas.Cmp(h.minPriorityFeePerGas) < 0 {
				priorityFeePerGas = h.minPriorityFeePerGas
			}
			if h.maxPriorityFeePerGas != nil && priorityFeePerGas.Cmp(h.maxPriorityFeePerGas) > 0 {
				priorityFeePerGas = h.maxPriorityFeePerGas
			}
			if maxFeePerGas.Cmp(priorityFeePerGas) < 0 {
				priorityFeePerGas = maxFeePerGas
			}
			dfd.MaxFeePerGas = maxFeePerGas
			dfd.MaxPriorityFeePerGas = priorityFeePerGas
			return next(ctx, t, result, method, args...)
		}
	}
}

func (h *hijackDynamicGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (h *hijackDynamicGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}
