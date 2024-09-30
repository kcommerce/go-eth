package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

type TransactionLegacy struct {
	EmbedCallData
	EmbedTransactionData
	EmbedLegacyPriceData
}

func NewTransactionLegacy() *TransactionLegacy {
	return &TransactionLegacy{}
}

func (t *TransactionLegacy) Type() TransactionType {
	return LegacyTxType
}

func (t *TransactionLegacy) Call() Call {
	return &CallLegacy{
		EmbedCallData:        *t.EmbedCallData.Copy(),
		EmbedLegacyPriceData: *t.EmbedLegacyPriceData.Copy(),
	}
}

func (t *TransactionLegacy) CalculateHash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

func (t *TransactionLegacy) CalculateSigningHash() (Hash, error) {
	var (
		chainID  = uint64(0)
		nonce    = uint64(0)
		gasPrice = big.NewInt(0)
		gasLimit = uint64(0)
		to       = ([]byte)(nil)
		value    = big.NewInt(0)
		input    = ([]byte)(nil)
	)
	if t.ChainID != nil {
		chainID = *t.ChainID
	}
	if t.Nonce != nil {
		nonce = *t.Nonce
	}
	if t.GasPrice != nil {
		gasPrice = t.GasPrice
	}
	if t.GasLimit != nil {
		gasLimit = *t.GasLimit
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = t.Value
	}
	if t.Input != nil {
		input = t.Input
	}
	list := rlp.List{
		rlp.Uint(nonce),
		(*rlp.BigInt)(gasPrice),
		rlp.Uint(gasLimit),
		rlp.Bytes(to),
		(*rlp.BigInt)(value),
		rlp.Bytes(input),
	}
	if t.ChainID != nil && *t.ChainID != 0 {
		list.Add(
			rlp.Uint(chainID),
			rlp.Uint(0),
			rlp.Uint(0),
		)
	}
	bin, err := list.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(bin)), nil
}

//nolint:funlen
func (t TransactionLegacy) EncodeRLP() ([]byte, error) {
	var (
		nonce    = uint64(0)
		gasPrice = big.NewInt(0)
		gasLimit = uint64(0)
		to       = ([]byte)(nil)
		value    = big.NewInt(0)
		input    = ([]byte)(nil)
		v        = big.NewInt(0)
		r        = big.NewInt(0)
		s        = big.NewInt(0)
	)
	if t.Nonce != nil {
		nonce = *t.Nonce
	}
	if t.GasPrice != nil {
		gasPrice = t.GasPrice
	}
	if t.GasLimit != nil {
		gasLimit = *t.GasLimit
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = t.Value
	}
	if t.Input != nil {
		input = t.Input
	}
	if t.Signature != nil {
		v = t.Signature.V
		r = t.Signature.R
		s = t.Signature.S
	}
	return rlp.List{
		rlp.Uint(nonce),
		(*rlp.BigInt)(gasPrice),
		rlp.Uint(gasLimit),
		rlp.Bytes(to),
		(*rlp.BigInt)(value),
		rlp.Bytes(input),
		(*rlp.BigInt)(v),
		(*rlp.BigInt)(r),
		(*rlp.BigInt)(s),
	}.EncodeRLP()
}

//nolint:funlen
func (t *TransactionLegacy) DecodeRLP(data []byte) (int, error) {
	*t = TransactionLegacy{}
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	var (
		nonce    = new(rlp.Uint)
		gasPrice = new(rlp.BigInt)
		gasLimit = new(rlp.Uint)
		to       = new(rlp.Bytes)
		value    = new(rlp.BigInt)
		input    = new(rlp.Bytes)
		v        = new(rlp.BigInt)
		r        = new(rlp.BigInt)
		s        = new(rlp.BigInt)
	)
	list := rlp.List{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		v,
		r,
		s,
	}
	if _, err := rlp.Decode(data, &list); err != nil {
		return 0, err
	}
	if nonce.Get() != 0 {
		t.Nonce = nonce.Ptr()
	}
	if gasPrice.Get().Sign() != 0 {
		t.GasPrice = gasPrice.Ptr()
	}
	if gasLimit.Get() != 0 {
		t.GasLimit = gasLimit.Ptr()
	}
	if len(to.Get()) > 0 {
		t.To = AddressFromBytesPtr(*to)
	}
	if value.Get().Sign() != 0 {
		t.Value = value.Ptr()
	}
	if len(input.Get()) > 0 {
		t.Input = input.Get()
	}
	if v.Get().Sign() != 0 || r.Get().Sign() != 0 || s.Get().Sign() != 0 {
		t.Signature = &Signature{
			V: (*big.Int)(v),
			R: (*big.Int)(r),
			S: (*big.Int)(s),
		}
		// Derive chain ID from the V value.
		if v.Get().Cmp(big.NewInt(35)) >= 0 {
			x := new(big.Int).Sub((*big.Int)(v), big.NewInt(35))
			x = x.Div(x, big.NewInt(2))
			chainID := x.Uint64()
			t.ChainID = &chainID
		}
	}
	return len(data), nil
}

func (t *TransactionLegacy) MarshalJSON() ([]byte, error) {
	transaction := &jsonTransactionLegacy{}
	transaction.To = t.To
	transaction.From = t.From
	if t.ChainID != nil {
		transaction.ChainID = NumberFromUint64Ptr(*t.ChainID)
	}
	if t.GasLimit != nil {
		transaction.GasLimit = NumberFromUint64Ptr(*t.GasLimit)
	}
	if t.GasPrice != nil {
		transaction.GasPrice = NumberFromBigIntPtr(t.GasPrice)
	}
	transaction.Input = t.Input
	if t.Nonce != nil {
		transaction.Nonce = NumberFromUint64Ptr(*t.Nonce)
	}
	if t.Value != nil {
		transaction.Value = NumberFromBigIntPtr(t.Value)
	}
	if t.Signature != nil {
		transaction.V = NumberFromBigIntPtr(t.Signature.V)
		transaction.R = NumberFromBigIntPtr(t.Signature.R)
		transaction.S = NumberFromBigIntPtr(t.Signature.S)
	}
	return json.Marshal(transaction)
}

func (t *TransactionLegacy) UnmarshalJSON(data []byte) error {
	transaction := &jsonTransactionLegacy{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return err
	}
	if transaction.ChainID != nil {
		chainID := transaction.ChainID.Big().Uint64()
		t.ChainID = &chainID
	}
	t.To = transaction.To
	t.From = transaction.From
	if transaction.GasLimit != nil {
		gas := transaction.GasLimit.Big().Uint64()
		t.GasLimit = &gas
	}
	if transaction.GasPrice != nil {
		t.GasPrice = transaction.GasPrice.Big()
	}
	t.Input = transaction.Input
	if transaction.Nonce != nil {
		nonce := transaction.Nonce.Big().Uint64()
		t.Nonce = &nonce
	}
	if transaction.Value != nil {
		t.Value = transaction.Value.Big()
	}
	if transaction.V != nil && transaction.R != nil && transaction.S != nil {
		t.Signature = SignatureFromVRSPtr(transaction.V.Big(), transaction.R.Big(), transaction.S.Big())
	}
	return nil
}

type jsonTransactionLegacy struct {
	ChainID  *Number  `json:"chainId,omitempty"`
	From     *Address `json:"from,omitempty"`
	To       *Address `json:"to,omitempty"`
	GasLimit *Number  `json:"gas,omitempty"`
	GasPrice *Number  `json:"gasPrice,omitempty"`
	Input    Bytes    `json:"input,omitempty"`
	Nonce    *Number  `json:"nonce,omitempty"`
	Value    *Number  `json:"value,omitempty"`
	V        *Number  `json:"v,omitempty"`
	R        *Number  `json:"r,omitempty"`
	S        *Number  `json:"s,omitempty"`
}

var _ Transaction = (*TransactionLegacy)(nil)
