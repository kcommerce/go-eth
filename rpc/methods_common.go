package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// MethodsCommon is a collection of methods that are commonly supported by
// Ethereum JSON-RPC APIs.
type MethodsCommon struct {
	Transport transport.Transport
	Decoder   types.TransactionDecoder
}

//
// Account methods:
//

// GetBalance performs eth_getBalance RPC call.
//
// It returns the balance of the account of given address in wei.
func (c *MethodsCommon) GetBalance(ctx context.Context, address types.Address, block types.BlockNumber) (*big.Int, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getBalance", address, block); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// GetCode performs eth_getCode RPC call.
//
// It returns the contract code at the given address.
func (c *MethodsCommon) GetCode(ctx context.Context, account types.Address, block types.BlockNumber) ([]byte, error) {
	var res types.Bytes
	if err := c.Transport.Call(ctx, &res, "eth_getCode", account, block); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

// GetStorageAt performs eth_getStorageAt RPC call.
//
// It returns the value of key in the contract storage at the given
// address.
func (c *MethodsCommon) GetStorageAt(ctx context.Context, account types.Address, key types.Hash, block types.BlockNumber) (*types.Hash, error) {
	var res types.Hash
	if err := c.Transport.Call(ctx, &res, "eth_getStorageAt", account, key, block); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionCount performs eth_getTransactionCount RPC call.
//
// It returns the number of transactions sent from the given address.
func (c *MethodsCommon) GetTransactionCount(ctx context.Context, account types.Address, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getTransactionCount", account, block); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

//
// Block methods:
//

// BlockByHash performs eth_getBlockByHash RPC call.
//
// It returns information about a block by hash.
func (c *MethodsCommon) BlockByHash(ctx context.Context, hash types.Hash, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.Transport.Call(ctx, &res, "eth_getBlockByHash", hash, full); err != nil {
		return nil, err
	}
	return &res, nil
}

// BlockByNumber performs eth_getBlockByNumber RPC call.
//
// It returns the block with the given number.
func (c *MethodsCommon) BlockByNumber(ctx context.Context, number types.BlockNumber, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.Transport.Call(ctx, &res, "eth_getBlockByNumber", number, full); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetBlockTransactionCountByHash performs eth_getBlockTransactionCountByHash RPC call.
//
// It returns the number of transactions in the block with the given hash.
func (c *MethodsCommon) GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getBlockTransactionCountByHash", hash); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetBlockTransactionCountByNumber implements the RPC interface.
func (c *MethodsCommon) GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getBlockTransactionCountByNumber", number); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetUncleByBlockHashAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
//
// It returns information about an uncle of a block by number and uncle index position.
func (c *MethodsCommon) GetUncleByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.Block, error) {
	var res types.Block
	if err := c.Transport.Call(ctx, &res, "eth_getUncleByBlockHashAndIndex", hash, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetUncleByBlockNumberAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
//
// It returns information about an uncle of a block by hash and uncle index position.
func (c *MethodsCommon) GetUncleByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.Block, error) {
	var res types.Block
	if err := c.Transport.Call(ctx, &res, "eth_getUncleByBlockNumberAndIndex", number, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetUncleCountByBlockHash performs eth_getUncleCountByBlockHash RPC call.
//
// It returns the number of uncles in the block with the given hash.
func (c *MethodsCommon) GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getUncleCountByBlockHash", hash); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("uncle count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetUncleCountByBlockNumber performs eth_getUncleCountByBlockNumber RPC call.
//
// It returns the number of uncles in the block with the given block number.
func (c *MethodsCommon) GetUncleCountByBlockNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_getUncleCountByBlockNumber", number); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("uncle count is too big")
	}
	return res.Big().Uint64(), nil
}

// SubscribeNewHeads performs eth_subscribe RPC call with "newHeads"
// subscription type.
//
// It creates a subscription that will send new block headers.
//
// Subscription channel will be closed when the context is canceled.
func (c *MethodsCommon) SubscribeNewHeads(ctx context.Context) (<-chan types.Block, error) {
	return subscribe[types.Block](ctx, c.Transport, "newHeads")
}

//
// Transaction methods:
//

// Call performs eth_call RPC call.
//
// It executes a new message call immediately without creating a
// transaction on the blockchain.
//
// If call also implements types.Transaction, then a Call method of the
// transaction will be used to create a call.
func (c *MethodsCommon) Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error) {
	if call == nil {
		return nil, errors.New("rpc client: call is nil")
	}
	if tx, ok := call.(types.Transaction); ok {
		call = tx.Call()
	}
	var res types.Bytes
	if err := c.Transport.Call(ctx, &res, "eth_call", call, block); err != nil {
		return nil, err
	}
	return res, nil
}

// EstimateGas performs eth_estimateGas RPC call.
//
// It estimates the gas necessary to execute a specific transaction.
//
// If call also implements types.Transaction, then a Call method of the
// transaction will be used to create a call.
func (c *MethodsCommon) EstimateGas(ctx context.Context, call types.Call, block types.BlockNumber) (uint64, error) {
	if call == nil {
		return 0, errors.New("rpc client: call is nil")
	}
	if tx, ok := call.(types.Transaction); ok {
		call = tx.Call()
	}
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_estimateGas", call, block); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("gas estimate is too big")
	}
	return res.Big().Uint64(), nil
}

// SendRawTransaction performs eth_sendRawTransaction RPC call.
//
// It sends an encoded transaction to the network.
func (c *MethodsCommon) SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error) {
	var res types.Hash
	if err := c.Transport.Call(ctx, &res, "eth_sendRawTransaction", types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByHash performs eth_getTransactionByHash RPC call.
//
// It returns the information about a transaction requested by transaction.
func (c *MethodsCommon) GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.TransactionOnChain, error) {
	res := types.TransactionOnChain{Decoder: c.Decoder}
	if err := c.Transport.Call(ctx, &res, "eth_getTransactionByHash", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByBlockHashAndIndex performs eth_getTransactionByBlockHashAndIndex RPC call.
//
// It returns the information about a transaction requested by transaction.
func (c *MethodsCommon) GetTransactionByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.TransactionOnChain, error) {
	res := types.TransactionOnChain{Decoder: c.Decoder}
	if err := c.Transport.Call(ctx, &res, "eth_getTransactionByBlockHashAndIndex", hash, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByBlockNumberAndIndex performs eth_getTransactionByBlockNumberAndIndex RPC call.
//
// It returns the information about a transaction requested by transaction.
func (c *MethodsCommon) GetTransactionByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.TransactionOnChain, error) {
	res := types.TransactionOnChain{Decoder: c.Decoder}
	if err := c.Transport.Call(ctx, &res, "eth_getTransactionByBlockNumberAndIndex", number, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionReceipt performs eth_getTransactionReceipt RPC call.
//
// It returns the receipt of a transaction by transaction hash.
func (c *MethodsCommon) GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	var res types.TransactionReceipt
	if err := c.Transport.Call(ctx, &res, "eth_getTransactionReceipt", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetBlockReceipts performs eth_getBlockReceipts RPC call.
//
// It returns all transaction receipts for a given block hash or number.
func (c *MethodsCommon) GetBlockReceipts(ctx context.Context, block types.BlockNumber) ([]*types.TransactionReceipt, error) {
	var res []*types.TransactionReceipt
	if err := c.Transport.Call(ctx, &res, "eth_getBlockReceipts", block); err != nil {
		return nil, err
	}
	return res, nil
}

// SubscribeNewPendingTransactions performs eth_subscribe RPC call with
// "newPendingTransactions" subscription type.
//
// It creates a subscription that will send new pending transactions.
//
// Subscription channel will be closed when the context is canceled.
func (c *MethodsCommon) SubscribeNewPendingTransactions(ctx context.Context) (<-chan types.Hash, error) {
	return subscribe[types.Hash](ctx, c.Transport, "newPendingTransactions")
}

//
// Logs methods:
//

// GetLogs performs eth_getLogs RPC call.
//
// It returns logs that match the given query.
func (c *MethodsCommon) GetLogs(ctx context.Context, query *types.FilterLogsQuery) ([]types.Log, error) {
	var res []types.Log
	if err := c.Transport.Call(ctx, &res, "eth_getLogs", query); err != nil {
		return nil, err
	}
	return res, nil
}

// SubscribeLogs performs eth_subscribe RPC call with "logs" subscription
// type.
//
// It creates a subscription that will send logs that match the given query.
//
// Subscription channel will be closed when the context is canceled.
func (c *MethodsCommon) SubscribeLogs(ctx context.Context, query *types.FilterLogsQuery) (<-chan types.Log, error) {
	return subscribe[types.Log](ctx, c.Transport, "logs", query)
}

// Network status methods:

// ChainID performs eth_chainId RPC call.
//
// It returns the current chain ID.
func (c *MethodsCommon) ChainID(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_chainId"); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, fmt.Errorf("chain id is too big")
	}
	return res.Big().Uint64(), nil
}

// BlockNumber performs eth_blockNumber RPC call.
//
// It returns the current block number.
func (c *MethodsCommon) BlockNumber(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_blockNumber"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// GasPrice performs eth_gasPrice RPC call.
//
// It returns the current price per gas in wei.
func (c *MethodsCommon) GasPrice(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// MaxPriorityFeePerGas performs eth_maxPriorityFeePerGas RPC call.
//
// It returns the estimated maximum priority fee per gas.
func (c *MethodsCommon) MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// BlobBaseFee performs eth_blobBaseFee RPC call.
//
// It returns the expected base fee for blobs in the next block.
func (c *MethodsCommon) BlobBaseFee(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.Transport.Call(ctx, &res, "eth_blobBaseFee"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}
