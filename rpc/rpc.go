package rpc

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// RPC is an RPC client for the Ethereum-compatible nodes.
type RPC2 interface {

	// GetBlockTransactionCountByHash performs eth_getBlockTransactionCountByHash RPC call.
	//
	// It returns the number of transactions in the block with the given hash.
	GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error)

	// GetBlockTransactionCountByNumber performs eth_getBlockTransactionCountByNumber RPC call.
	//
	// It returns the number of transactions in the block with the given block
	GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error)

	// GetUncleCountByBlockHash performs eth_getUncleCountByBlockHash RPC call.
	//
	// It returns the number of uncles in the block with the given hash.
	GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error)

	// Sign performs eth_sign RPC call.
	//
	// It signs the given data with the given address.
	Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error)

	// SignTransaction performs eth_signTransaction RPC call.
	//
	// It signs the given transaction.
	//
	// If transaction was internally mutated, the mutated call is returned.
	SignTransaction(ctx context.Context, tx types.Transaction) (types.Transaction, error)

	// SendTransaction performs eth_sendTransaction RPC call.
	//
	// It sends a transaction to the network.
	//
	// If transaction was internally mutated, the mutated call is returned.
	SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error)

	// SendRawTransaction performs eth_sendRawTransaction RPC call.
	//
	// It sends an encoded transaction to the network.
	SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error)

	// Call performs eth_call RPC call.
	//
	// It executes a new message call immediately without creating a
	// transaction on the blockchain.
	//
	// If call was internally mutated, the mutated call is returned.
	Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error)

	// EstimateGas performs eth_estimateGas RPC call.
	//
	// It estimates the gas necessary to execute a specific transaction.
	//
	// If call was internally mutated, the mutated call is returned.
	EstimateGas(ctx context.Context, call types.Call, block types.BlockNumber) (uint64, error)

	// BlockByHash performs eth_getBlockByHash RPC call.
	//
	// It returns information about a block by hash.
	BlockByHash(ctx context.Context, hash types.Hash, full bool) (*types.Block, error)

	// BlockByNumber performs eth_getBlockByNumber RPC call.
	//
	// It returns the block with the given number.
	BlockByNumber(ctx context.Context, number types.BlockNumber, full bool) (*types.Block, error)

	// GetTransactionByHash performs eth_getTransactionByHash RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.TransactionOnChain, error)

	// GetTransactionByBlockHashAndIndex performs eth_getTransactionByBlockHashAndIndex RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.TransactionOnChain, error)

	// GetTransactionByBlockNumberAndIndex performs eth_getTransactionByBlockNumberAndIndex RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.TransactionOnChain, error)

	// GetTransactionReceipt performs eth_getTransactionReceipt RPC call.
	//
	// It returns the receipt of a transaction by transaction hash.
	GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error)

	// GetBlockReceipts performs eth_getBlockReceipts RPC call.
	//
	// It returns all transaction receipts for a given block hash or number.
	GetBlockReceipts(ctx context.Context, block types.BlockNumber) ([]*types.TransactionReceipt, error)

	// GetUncleByBlockHashAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
	//
	// It returns information about an uncle of a block by number and uncle index position.
	GetUncleByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.Block, error)

	// GetUncleByBlockNumberAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
	//
	// It returns information about an uncle of a block by hash and uncle index position.
	GetUncleByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.Block, error)

	// NewFilter performs eth_newFilter RPC call.
	//
	// It creates a filter object based on the given filter options. To check
	// if the state has changed, use GetFilterChanges.
	NewFilter(ctx context.Context, query *types.FilterLogsQuery) (*big.Int, error)

	// NewBlockFilter performs eth_newBlockFilter RPC call.
	//
	// It creates a filter in the node, to notify when a new block arrives. To
	// check if the state has changed, use GetBlockFilterChanges.
	NewBlockFilter(ctx context.Context) (*big.Int, error)

	// NewPendingTransactionFilter performs eth_newPendingTransactionFilter RPC call.
	//
	// It creates a filter in the node, to notify when new pending transactions
	// arrive. To check if the state has changed, use GetFilterChanges.
	NewPendingTransactionFilter(ctx context.Context) (*big.Int, error)

	// UninstallFilter performs eth_uninstallFilter RPC call.
	//
	// It uninstalls a filter with given ID. Should always be called when watch
	// is no longer needed.
	UninstallFilter(ctx context.Context, id *big.Int) (bool, error)

	// GetFilterChanges performs eth_getFilterChanges RPC call.
	//
	// It returns an array of logs that occurred since the given filter ID.
	GetFilterChanges(ctx context.Context, id *big.Int) ([]types.Log, error)

	// GetBlockFilterChanges performs eth_getFilterChanges RPC call.
	//
	// It returns an array of block hashes that occurred since the given filter ID.
	GetBlockFilterChanges(ctx context.Context, id *big.Int) ([]types.Hash, error)

	// GetFilterLogs performs eth_getFilterLogs RPC call.
	//
	// It returns an array of all logs matching filter with given ID.
	GetFilterLogs(ctx context.Context, id *big.Int) ([]types.Log, error)

	// GetLogs performs eth_getLogs RPC call.
	//
	// It returns logs that match the given query.
	GetLogs(ctx context.Context, query *types.FilterLogsQuery) ([]types.Log, error)

	// MaxPriorityFeePerGas performs eth_maxPriorityFeePerGas RPC call.
	//
	// It returns the estimated maximum priority fee per gas.
	MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error)

	// SubscribeLogs performs eth_subscribe RPC call with "logs" subscription
	// type.
	//
	// It creates a subscription that will send logs that match the given query.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeLogs(ctx context.Context, query *types.FilterLogsQuery) (<-chan types.Log, error)

	// SubscribeNewHeads performs eth_subscribe RPC call with "newHeads"
	// subscription type.
	//
	// It creates a subscription that will send new block headers.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeNewHeads(ctx context.Context) (<-chan types.Block, error)

	// SubscribeNewPendingTransactions performs eth_subscribe RPC call with
	// "newPendingTransactions" subscription type.
	//
	// It creates a subscription that will send new pending transactions.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeNewPendingTransactions(ctx context.Context) (<-chan types.Hash, error)
}
