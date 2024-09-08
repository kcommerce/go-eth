package crypto

import (
	"github.com/defiweb/go-eth/crypto/ecdsa"
	"github.com/defiweb/go-eth/crypto/keccak"
	"github.com/defiweb/go-eth/crypto/kzg4844"
)

// Default implementations of the crypto functions.
var (
	Keccak256            = keccak.Hash256
	ECPublicKeyToAddress = ecdsa.PublicKeyToAddress
	ECSignHash           = ecdsa.SignHash
	ECRecoverHash        = ecdsa.RecoverHash
	ECSignMessage        = ecdsa.SignMessage
	ECRecoverMessage     = ecdsa.RecoverMessage
	KZGBlobToCommitment  = kzg4844.BlobToCommitment
	KZGComputeProof      = kzg4844.ComputeProof
	KZGVerifyProof       = kzg4844.VerifyProof
	KZGComputeBlobProof  = kzg4844.ComputeBlobProof
	KZGVerifyBlobProof   = kzg4844.VerifyBlobProof
	KZGComputeBlobHashV1 = kzg4844.ComputeBlobHashV1
)
