package types

import (
	"encoding/json"

	"github.com/defiweb/go-eth/crypto/kzg4844"
)

type CallBlob struct {
	EmbedCallData
	EmbedAccessListData
	EmbedDynamicFeeData
	EmbedBlobData
}

func NewCallBlob() *CallBlob {
	return &CallBlob{}
}

func (c *CallBlob) Copy() *CallBlob {
	if c == nil {
		return nil
	}
	return &CallBlob{
		EmbedCallData:       *c.EmbedCallData.Copy(),
		EmbedAccessListData: *c.EmbedAccessListData.Copy(),
		EmbedDynamicFeeData: *c.EmbedDynamicFeeData.Copy(),
		EmbedBlobData:       *c.EmbedBlobData.Copy(),
	}
}

func (c *CallBlob) MarshalJSON() ([]byte, error) {
	call := &jsonCallBlob{
		From:       c.From,
		To:         c.To,
		Data:       c.Input,
		AccessList: c.AccessList,
	}
	if c.GasLimit != nil {
		call.GasLimit = NumberFromUint64Ptr(*c.GasLimit)
	}
	if c.MaxFeePerGas != nil {
		call.MaxFeePerGas = NumberFromBigIntPtr(c.MaxFeePerGas)
	}
	if c.MaxFeePerBlobGas != nil {
		call.MaxFeePerBlobGas = NumberFromBigIntPtr(c.MaxFeePerBlobGas)
	}
	if c.MaxPriorityFeePerGas != nil {
		call.MaxPriorityFeePerGas = NumberFromBigIntPtr(c.MaxPriorityFeePerGas)
	}
	if c.Value != nil {
		value := NumberFromBigInt(c.Value)
		call.Value = &value
	}
	for _, blob := range c.Blobs {
		hash := blob.Hash
		if hash.IsZero() && blob.Sidecar != nil {
			hash = blob.Sidecar.ComputeHash()
		}
		call.BlobHashes = append(call.BlobHashes, hash)
		if blob.Sidecar != nil {
			call.Blobs = append(call.Blobs, kzgBlob(blob.Sidecar.Blob))
			call.Commitments = append(call.Commitments, kzgCommitment(blob.Sidecar.Commitment))
			call.Proofs = append(call.Proofs, kzgProof(blob.Sidecar.Proof))
		}
	}
	return json.Marshal(call)
}

func (c *CallBlob) UnmarshalJSON(data []byte) error {
	call := &jsonCallBlob{}
	if err := json.Unmarshal(data, call); err != nil {
		return err
	}
	c.From = call.From
	c.To = call.To
	if call.GasLimit != nil {
		gas := call.GasLimit.Big().Uint64()
		c.GasLimit = &gas
	}
	if call.MaxFeePerGas != nil {
		c.MaxFeePerGas = call.MaxFeePerGas.Big()
	}
	if call.MaxFeePerBlobGas != nil {
		c.MaxFeePerBlobGas = call.MaxFeePerBlobGas.Big()
	}
	if call.MaxPriorityFeePerGas != nil {
		c.MaxPriorityFeePerGas = call.MaxPriorityFeePerGas.Big()
	}
	if call.Value != nil {
		c.Value = call.Value.Big()
	}
	c.Input = call.Data
	c.AccessList = call.AccessList
	if len(call.BlobHashes) > 0 {
		c.Blobs = make([]Blob, len(call.BlobHashes))
		for i, hash := range call.BlobHashes {
			blob := Blob{Hash: hash}
			if i < len(call.Blobs) && i < len(call.Commitments) && i < len(call.Proofs) {
				blob.Sidecar = &BlobSidecar{
					Blob:       kzg4844.Blob(call.Blobs[i]),
					Commitment: kzg4844.Commitment(call.Commitments[i]),
					Proof:      kzg4844.Proof(call.Proofs[i]),
				}
			}
			c.Blobs[i] = blob
		}
	}
	return nil
}

type jsonCallBlob struct {
	From                 *Address        `json:"from,omitempty"`
	To                   *Address        `json:"to,omitempty"`
	GasLimit             *Number         `json:"gas,omitempty"`
	MaxFeePerGas         *Number         `json:"maxFeePerGas,omitempty"`
	MaxFeePerBlobGas     *Number         `json:"maxFeePerBlobGas,omitempty"`
	MaxPriorityFeePerGas *Number         `json:"maxPriorityFeePerGas,omitempty"`
	Value                *Number         `json:"value,omitempty"`
	Data                 Bytes           `json:"data,omitempty"`
	AccessList           AccessList      `json:"accessList,omitempty"`
	BlobHashes           []Hash          `json:"blobVersionedHashes,omitempty"`
	Blobs                []kzgBlob       `json:"blobs,omitempty"`
	Commitments          []kzgCommitment `json:"commitments,omitempty"`
	Proofs               []kzgProof      `json:"proofs,omitempty"`
}
