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
		input                = ([]byte)(nil)
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
	if t.Input != nil {
		input = t.Input
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
	bin, err := rlp.List{
		rlp.Uint(chainID),
		rlp.Uint(nonce),
		(*rlp.BigInt)(maxPriorityFeePerGas),
		(*rlp.BigInt)(maxFeePerGas),
		rlp.Uint(gasLimit),
		rlp.Bytes(to),
		(*rlp.BigInt)(value),
		rlp.Bytes(input),
		&accessList,
		(*rlp.BigInt)(maxFeePerBlobGas),
		&blobHashes,
	}.EncodeRLP()
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
		input                = ([]byte)(nil)
		accessList           = (AccessList)(nil)
		maxFeePerBlobGas     = big.NewInt(0)
		blobHashes           = (hashList)(nil)
		blobs                = rlp.TypedList[kzgBlob]{}
		commitments          = rlp.TypedList[kzgCommitment]{}
		proofs               = rlp.TypedList[kzgProof]{}
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
	if t.Input != nil {
		input = t.Input
	}
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	if t.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = t.MaxFeePerBlobGas
	}
	if len(t.Blobs) > 0 {
		blobHashes = make(hashList, 0, len(t.Blobs))
		for _, blob := range t.Blobs {
			hash := blob.Hash
			if hash.IsZero() && blob.Sidecar != nil {
				hash = blob.Sidecar.ComputeHash()
			}
			blobHashes = append(blobHashes, hash)
			if blob.Sidecar != nil {
				blobs.Add((*kzgBlob)(&blob.Sidecar.Blob))
				commitments.Add((*kzgCommitment)(&blob.Sidecar.Commitment))
				proofs.Add((*kzgProof)(&blob.Sidecar.Proof))
			}
		}
	}
	if t.Signature != nil {
		v = t.Signature.V
		r = t.Signature.R
		s = t.Signature.S
	}
	tx := rlp.List{
		rlp.Uint(chainID),
		rlp.Uint(nonce),
		(*rlp.BigInt)(maxPriorityFeePerGas),
		(*rlp.BigInt)(maxFeePerGas),
		rlp.Uint(gasLimit),
		rlp.Bytes(to),
		(*rlp.BigInt)(value),
		rlp.Bytes(input),
		&accessList,
		(*rlp.BigInt)(maxFeePerBlobGas),
		&blobHashes,
		(*rlp.BigInt)(v),
		(*rlp.BigInt)(r),
		(*rlp.BigInt)(s),
	}
	if len(blobHashes) > 0 && len(blobHashes) == len(blobs) {
		tx = rlp.List{
			tx,
			blobs,
			commitments,
			proofs,
		}
	}
	bin, err := tx.EncodeRLP()
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
		chainID              = new(rlp.Uint)
		nonce                = new(rlp.Uint)
		gasLimit             = new(rlp.Uint)
		maxPriorityFeePerGas = new(rlp.BigInt)
		maxFeePerGas         = new(rlp.BigInt)
		to                   = new(rlp.Bytes)
		value                = new(rlp.BigInt)
		input                = new(rlp.Bytes)
		accessList           = new(AccessList)
		maxFeePerBlobGas     = new(rlp.BigInt)
		blobHashes           = &hashList{}
		blobs                = new(rlp.TypedList[kzgBlob])
		commitments          = new(rlp.TypedList[kzgCommitment])
		proofs               = new(rlp.TypedList[kzgProof])
		v                    = new(rlp.BigInt)
		r                    = new(rlp.BigInt)
		s                    = new(rlp.BigInt)
	)
	dec, _, err := rlp.DecodeLazy(data)
	if err != nil {
		return 0, err
	}
	if !dec.IsList() {
		return 0, fmt.Errorf("unable to decode transaction")
	}
	var list rlp.List
	switch dec.Length() {
	case 4:
		list = rlp.List{
			&rlp.List{
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
			},
			blobs,
			commitments,
			proofs,
		}
	default:
		list = rlp.List{
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
		}
	}
	if err := dec.Decode(&list); err != nil {
		return 0, err
	}
	if chainID.Get() != 0 {
		t.ChainID = chainID.Ptr()
	}
	if nonce.Get() != 0 {
		t.Nonce = nonce.Ptr()
	}
	if maxPriorityFeePerGas.Ptr().Sign() != 0 {
		t.MaxPriorityFeePerGas = maxPriorityFeePerGas.Ptr()
	}
	if maxFeePerGas.Ptr().Sign() != 0 {
		t.MaxFeePerGas = maxFeePerGas.Ptr()
	}
	if gasLimit.Get() != 0 {
		t.GasLimit = gasLimit.Ptr()
	}
	if len(to.Get()) > 0 {
		t.To = AddressFromBytesPtr(to.Get())
	}
	if value.Ptr().Sign() != 0 {
		t.Value = value.Ptr()
	}
	if len(input.Get()) > 0 {
		t.Input = input.Get()
	}
	if len(*accessList) > 0 {
		t.AccessList = *accessList
	}
	if maxFeePerBlobGas.Ptr().Sign() != 0 {
		t.MaxFeePerBlobGas = maxFeePerBlobGas.Ptr()
	}
	if len(*blobHashes) > 0 {
		t.Blobs = make([]Blob, len(*blobHashes))
		for i, hash := range *blobHashes {
			blob := Blob{Hash: hash}
			if i < len(*blobs) && i < len(*commitments) && i < len(*proofs) {
				blob.Sidecar = &BlobSidecar{
					Blob:       kzg4844.Blob(*(*blobs)[i]),
					Commitment: kzg4844.Commitment(*(*commitments)[i]),
					Proof:      kzg4844.Proof(*(*proofs)[i]),
				}
			}
			t.Blobs[i] = blob
		}
	}
	if v.Ptr().Sign() != 0 || r.Ptr().Sign() != 0 || s.Ptr().Sign() != 0 {
		t.Signature = &Signature{
			V: v.Ptr(),
			R: r.Ptr(),
			S: s.Ptr(),
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
