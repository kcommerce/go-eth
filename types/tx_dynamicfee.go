package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/defiweb/go-rlp"
)

type TransactionDynamicFee struct {
	EmbedCallData
	EmbedTransactionData
	EmbedAccessListData
	EmbedDynamicFeeData
}

func NewTransactionDynamicFee() *TransactionDynamicFee {
	return &TransactionDynamicFee{}
}

func (t *TransactionDynamicFee) Type() TransactionType {
	return DynamicFeeTxType
}

func (t *TransactionDynamicFee) Call() Call {
	return &CallDynamicFee{
		EmbedCallData:       *t.EmbedCallData.Copy(),
		EmbedAccessListData: *t.EmbedAccessListData.Copy(),
		EmbedDynamicFeeData: *t.EmbedDynamicFeeData.Copy(),
	}
}

func (t *TransactionDynamicFee) CalculateHash(h HashFunc) (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return h(raw), nil
}

func (t *TransactionDynamicFee) CalculateSigningHash(h HashFunc) (Hash, error) {
	var (
		chainID              = uint64(0)
		nonce                = uint64(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (AccessList)(nil)
	)
	if t.ChainID != nil {
		chainID = *t.ChainID
	}
	if t.Nonce != nil {
		nonce = *t.Nonce
	}
	if t.GasLimit != nil {
		gasLimit = *t.GasLimit
	}
	if t.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = t.MaxPriorityFeePerGas
	}
	if t.MaxFeePerGas != nil {
		maxFeePerGas = t.MaxFeePerGas
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = t.Value
	}
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	bin, err := rlp.NewList(
		rlp.NewUint(chainID),
		rlp.NewUint(nonce),
		rlp.NewBigInt(maxPriorityFeePerGas),
		rlp.NewBigInt(maxFeePerGas),
		rlp.NewUint(gasLimit),
		rlp.NewBytes(to),
		rlp.NewBigInt(value),
		rlp.NewBytes(t.Input),
		&accessList,
	).EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	bin = append([]byte{byte(DynamicFeeTxType)}, bin...)
	return h(bin), nil
}

//nolint:funlen
func (t TransactionDynamicFee) EncodeRLP() ([]byte, error) {
	var (
		chainID              = uint64(0)
		nonce                = uint64(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (AccessList)(nil)
		v                    = big.NewInt(0)
		r                    = big.NewInt(0)
		s                    = big.NewInt(0)
	)
	if t.ChainID != nil {
		chainID = *t.ChainID
	}
	if t.Nonce != nil {
		nonce = *t.Nonce
	}
	if t.GasLimit != nil {
		gasLimit = *t.GasLimit
	}
	if t.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = t.MaxPriorityFeePerGas
	}
	if t.MaxFeePerGas != nil {
		maxFeePerGas = t.MaxFeePerGas
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = t.Value
	}
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	if t.Signature != nil {
		v = t.Signature.V
		r = t.Signature.R
		s = t.Signature.S
	}
	bin, err := rlp.NewList(
		rlp.NewUint(chainID),
		rlp.NewUint(nonce),
		rlp.NewBigInt(maxPriorityFeePerGas),
		rlp.NewBigInt(maxFeePerGas),
		rlp.NewUint(gasLimit),
		rlp.NewBytes(to),
		rlp.NewBigInt(value),
		rlp.NewBytes(t.Input),
		&accessList,
		rlp.NewBigInt(v),
		rlp.NewBigInt(r),
		rlp.NewBigInt(s),
	).EncodeRLP()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(DynamicFeeTxType)}, bin...), nil
}

//nolint:funlen
func (t *TransactionDynamicFee) DecodeRLP(data []byte) (int, error) {
	*t = TransactionDynamicFee{}
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(DynamicFeeTxType) {
		return 0, fmt.Errorf("invalid transaction type: %d", data[0])
	}
	data = data[1:]
	var (
		list                 *rlp.ListItem
		chainID              = &rlp.UintItem{}
		nonce                = &rlp.UintItem{}
		gasLimit             = &rlp.UintItem{}
		maxPriorityFeePerGas = &rlp.BigIntItem{}
		maxFeePerGas         = &rlp.BigIntItem{}
		to                   = &rlp.StringItem{}
		value                = &rlp.BigIntItem{}
		input                = &rlp.StringItem{}
		accessList           = &AccessList{}
		v                    = &rlp.BigIntItem{}
		r                    = &rlp.BigIntItem{}
		s                    = &rlp.BigIntItem{}
	)
	list = rlp.NewList(
		chainID,
		nonce,
		maxPriorityFeePerGas,
		maxFeePerGas,
		gasLimit,
		to,
		value,
		input,
		accessList,
		v,
		r,
		s,
	)
	if _, err := rlp.DecodeTo(data, list); err != nil {
		return 0, err
	}
	if chainID.X != 0 {
		t.ChainID = &chainID.X
	}
	if nonce.X != 0 {
		t.Nonce = &nonce.X
	}
	if maxPriorityFeePerGas.X.Sign() != 0 {
		t.MaxPriorityFeePerGas = maxPriorityFeePerGas.X
	}
	if maxFeePerGas.X.Sign() != 0 {
		t.MaxFeePerGas = maxFeePerGas.X
	}
	if gasLimit.X != 0 {
		t.GasLimit = &gasLimit.X
	}
	if len(to.Bytes()) > 0 {
		t.To = AddressFromBytesPtr(to.Bytes())
	}
	if value.X.Sign() != 0 {
		t.Value = value.X
	}
	if len(input.Bytes()) > 0 {
		t.Input = input.Bytes()
	}
	if len(*accessList) > 0 {
		t.AccessList = *accessList
	}
	if v.X.Sign() != 0 || r.X.Sign() != 0 || s.X.Sign() != 0 {
		t.Signature = &Signature{
			V: v.X,
			R: r.X,
			S: s.X,
		}
	}
	return len(data), nil
}

func (t *TransactionDynamicFee) MarshalJSON() ([]byte, error) {
	transaction := &jsonTransactionDynamicFee{}
	if t.ChainID != nil {
		transaction.ChainID = NumberFromUint64Ptr(*t.ChainID)
	}
	transaction.To = t.To
	transaction.From = t.From
	if t.GasLimit != nil {
		transaction.GasLimit = NumberFromUint64Ptr(*t.GasLimit)
	}
	if t.MaxFeePerGas != nil {
		transaction.MaxFeePerGas = NumberFromBigIntPtr(t.MaxFeePerGas)
	}
	if t.MaxPriorityFeePerGas != nil {
		transaction.MaxPriorityFeePerGas = NumberFromBigIntPtr(t.MaxPriorityFeePerGas)
	}
	transaction.Input = t.Input
	if t.Nonce != nil {
		transaction.Nonce = NumberFromUint64Ptr(*t.Nonce)
	}
	if t.Value != nil {
		transaction.Value = NumberFromBigIntPtr(t.Value)
	}
	transaction.AccessList = t.AccessList
	if t.Signature != nil {
		transaction.V = NumberFromBigIntPtr(t.Signature.V)
		transaction.R = NumberFromBigIntPtr(t.Signature.R)
		transaction.S = NumberFromBigIntPtr(t.Signature.S)
	}
	return json.Marshal(transaction)
}

func (t *TransactionDynamicFee) UnmarshalJSON(data []byte) error {
	transaction := &jsonTransactionDynamicFee{}
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
	if transaction.MaxFeePerGas != nil {
		t.MaxFeePerGas = transaction.MaxFeePerGas.Big()
	}
	if transaction.MaxPriorityFeePerGas != nil {
		t.MaxPriorityFeePerGas = transaction.MaxPriorityFeePerGas.Big()
	}
	t.Input = transaction.Input
	if transaction.Nonce != nil {
		Nonce := transaction.Nonce.Big().Uint64()
		t.Nonce = &Nonce
	}
	if transaction.Value != nil {
		t.Value = transaction.Value.Big()
	}
	t.AccessList = transaction.AccessList
	if transaction.V != nil && transaction.R != nil && transaction.S != nil {
		t.Signature = SignatureFromVRSPtr(transaction.V.Big(), transaction.R.Big(), transaction.S.Big())
	}
	return nil
}

type jsonTransactionDynamicFee struct {
	ChainID              *Number    `json:"chainId,omitempty"`
	From                 *Address   `json:"from,omitempty"`
	To                   *Address   `json:"to,omitempty"`
	GasLimit             *Number    `json:"gas,omitempty"`
	MaxFeePerGas         *Number    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *Number    `json:"maxPriorityFeePerGas,omitempty"`
	Input                Bytes      `json:"input,omitempty"`
	Nonce                *Number    `json:"nonce,omitempty"`
	Value                *Number    `json:"value,omitempty"`
	AccessList           AccessList `json:"accessList,omitempty"`
	V                    *Number    `json:"v,omitempty"`
	R                    *Number    `json:"r,omitempty"`
	S                    *Number    `json:"s,omitempty"`
}

var _ Transaction = (*TransactionDynamicFee)(nil)
