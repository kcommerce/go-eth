package types

import (
	"encoding/json"

	"github.com/defiweb/go-rlp"
)

// TransactionType is the type of transaction.
type TransactionType uint8

// Transaction types.
const (
	LegacyTxType TransactionType = iota
	AccessListTxType
	DynamicFeeTxType
	BlobTxType
)

type Transaction interface {
	json.Marshaler
	json.Unmarshaler
	rlp.Encoder
	rlp.Decoder

	HasTransactionData

	// Type returns the type of the transaction.
	Type() TransactionType

	// Call returns the call associated with the transaction. The call is a
	// copy and can be modified. It may return nil if it is impossible to
	// create a call.
	Call() Call

	// CalculateHash calculates the hash of the transaction.
	CalculateHash() (Hash, error)

	// CalculateSigningHash calculates the signing hash of the transaction.
	CalculateSigningHash() (Hash, error)
}

// TransactionDecoder is an interface that is used to decode transactions of
// unknown types.
//
// Decoder may not set the From field of the transaction.
// To get signer of the transaction, use the Recoverer interface.
type TransactionDecoder interface {
	RPCTransactionDecoder
	JSONTransactionDecoder
}

// RPCTransactionDecoder is an interface that is used to decode transactions
// from RLP encoded data.
type RPCTransactionDecoder interface {
	// DecodeRLP decodes the RLP encoded transaction data.
	DecodeRLP(data []byte) (Transaction, error)
}

// JSONTransactionDecoder is an interface that is used to decode transactions
// from JSON encoded data.
type JSONTransactionDecoder interface {
	// DecodeJSON decodes the JSON encoded transaction data.
	DecodeJSON(data []byte) (Transaction, error)
}
