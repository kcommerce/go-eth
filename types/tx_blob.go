package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/kzg4844"
)

type TransactionBlob struct {
	EmbedCallData
	EmbedTransactionData
	EmbedAccessListData
	EmbedDynamicFeeData
	EmbedBlobData
}

func NewTransactionBlob() *TransactionBlob {
	return &TransactionBlob{}
}

func (t *TransactionBlob) Type() TransactionType {
	return BlobTxType
}

func (t *TransactionBlob) Call() Call {
	return &CallBlob{
		EmbedCallData:       *t.EmbedCallData.Copy(),
		EmbedAccessListData: *t.EmbedAccessListData.Copy(),
		EmbedDynamicFeeData: *t.EmbedDynamicFeeData.Copy(),
		EmbedBlobData:       *t.EmbedBlobData.Copy(),
	}
}

func (t *TransactionBlob) CalculateHash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

func (t *TransactionBlob) CalculateSigningHash() (Hash, error) {
	var (
		chainID              = uint64(0)
		nonce                = uint64(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (AccessList)(nil)
		maxFeePerBlobGas     = big.NewInt(0)
		blobHashes           = (hashList)(nil)
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
	if t.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = t.MaxFeePerBlobGas
	}
	if len(t.Blobs) > 0 {
		blobHashes = make(hashList, len(t.Blobs))
		for i, blob := range t.Blobs {
			if blob.Hash.IsZero() && blob.Sidecar != nil {
				blobHashes[i] = blob.Sidecar.ComputeHash()
				continue
			}
			blobHashes[i] = blob.Hash
		}
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
		rlp.NewBigInt(maxFeePerBlobGas),
		&blobHashes,
	).EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	bin = append([]byte{byte(BlobTxType)}, bin...)
	return Hash(crypto.Keccak256(bin)), nil
}

//nolint:funlen
func (t TransactionBlob) EncodeRLP() ([]byte, error) {
	var (
		chainID              = uint64(0)
		nonce                = uint64(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (AccessList)(nil)
		maxFeePerBlobGas     = big.NewInt(0)
		blobHashes           = (hashList)(nil)
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
	if t.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = t.MaxFeePerBlobGas
	}
	if len(t.Blobs) > 0 {
		blobHashes = make(hashList, len(t.Blobs))
		for i, blob := range t.Blobs {
			if blob.Hash.IsZero() && blob.Sidecar != nil {
				blobHashes[i] = blob.Sidecar.ComputeHash()
				continue
			}
			blobHashes[i] = blob.Hash
		}
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
		rlp.NewBigInt(maxFeePerBlobGas),
		&blobHashes,
		rlp.NewBigInt(v),
		rlp.NewBigInt(r),
		rlp.NewBigInt(s),
	).EncodeRLP()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(BlobTxType)}, bin...), nil
}

//nolint:funlen
func (t *TransactionBlob) DecodeRLP(data []byte) (int, error) {
	*t = TransactionBlob{}
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(BlobTxType) {
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
		maxFeePerBlobGas     = &rlp.BigIntItem{}
		blobHashes           = &hashList{}
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
		maxFeePerBlobGas,
		blobHashes,
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
	if maxFeePerBlobGas.X.Sign() != 0 {
		t.MaxFeePerBlobGas = maxFeePerBlobGas.X
	}
	if len(*blobHashes) > 0 {
		t.Blobs = make([]Blob, len(*blobHashes))
		for i, hash := range *blobHashes {
			t.Blobs[i] = Blob{Hash: hash}
		}
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

func (t *TransactionBlob) MarshalJSON() ([]byte, error) {
	transaction := &jsonTransactionBlob{}
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
	if t.MaxFeePerBlobGas != nil {
		transaction.MaxFeePerBlobGas = NumberFromBigIntPtr(t.MaxFeePerBlobGas)
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
	for _, blob := range t.Blobs {
		hash := blob.Hash
		if hash.IsZero() && blob.Sidecar != nil {
			hash = blob.Sidecar.ComputeHash()
		}
		transaction.BlobHashes = append(transaction.BlobHashes, hash)
		if blob.Sidecar != nil {
			transaction.Blobs = append(transaction.Blobs, kzgBlob(blob.Sidecar.Blob))
			transaction.Commitments = append(transaction.Commitments, kzgCommitment(blob.Sidecar.Commitment))
			transaction.Proofs = append(transaction.Proofs, kzgProof(blob.Sidecar.Proof))
		}
	}
	return json.Marshal(transaction)
}

func (t *TransactionBlob) UnmarshalJSON(data []byte) error {
	transaction := &jsonTransactionBlob{}
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
	if transaction.MaxFeePerBlobGas != nil {
		t.MaxFeePerBlobGas = transaction.MaxFeePerBlobGas.Big()
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
	if len(transaction.BlobHashes) > 0 {
		t.Blobs = make([]Blob, len(transaction.BlobHashes))
		for i, hash := range transaction.BlobHashes {
			blob := Blob{Hash: hash}
			if i < len(transaction.Blobs) && i < len(transaction.Commitments) && i < len(transaction.Proofs) {
				blob.Sidecar = &BlobSidecar{
					Blob:       kzg4844.Blob(transaction.Blobs[i]),
					Commitment: kzg4844.Commitment(transaction.Commitments[i]),
					Proof:      kzg4844.Proof(transaction.Proofs[i]),
				}
			}
			t.Blobs[i] = blob
		}
	}
	return nil
}

type jsonTransactionBlob struct {
	ChainID              *Number         `json:"chainId,omitempty"`
	From                 *Address        `json:"from,omitempty"`
	To                   *Address        `json:"to,omitempty"`
	GasLimit             *Number         `json:"gas,omitempty"`
	MaxFeePerGas         *Number         `json:"maxFeePerGas,omitempty"`
	MaxFeePerBlobGas     *Number         `json:"maxFeePerBlobGas,omitempty"`
	MaxPriorityFeePerGas *Number         `json:"maxPriorityFeePerGas,omitempty"`
	Input                Bytes           `json:"input,omitempty"`
	Nonce                *Number         `json:"nonce,omitempty"`
	Value                *Number         `json:"value,omitempty"`
	AccessList           AccessList      `json:"accessList,omitempty"`
	BlobHashes           []Hash          `json:"blobVersionedHashes,omitempty"`
	Blobs                []kzgBlob       `json:"blobs,omitempty"`
	Commitments          []kzgCommitment `json:"commitments,omitempty"`
	Proofs               []kzgProof      `json:"proofs,omitempty"`
	V                    *Number         `json:"v,omitempty"`
	R                    *Number         `json:"r,omitempty"`
	S                    *Number         `json:"s,omitempty"`
}

var _ Transaction = (*TransactionBlob)(nil)
