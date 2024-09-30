package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	btcececdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"

	"github.com/defiweb/go-eth/crypto/keccak"
)

// PublicKey is an ECDSA public key.
type PublicKey = ecdsa.PublicKey

// PrivateKey is an ECDSA private key.
type PrivateKey = ecdsa.PrivateKey

// Hash is a 32-byte hash.
//
// This type is used in ECDSA cryptographic operations. For most use cases, the
// Hash type from the types package should be used instead.
type Hash [32]byte

// Address is a 20-byte Ethereum address.
//
// This type is used in ECDSA cryptographic operations. For most use cases, the
// Address type from the types package should be used instead.
type Address [20]byte

// Signature is an ECDSA signature.
//
// This type is used in ECDSA cryptographic operations. For most use cases, the
// Signature type from the types package should be used instead.
type Signature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}

var s256 = btcec.S256()

// AddMessagePrefix adds the Ethereum message prefix to the given data as
// defined in EIP-191.
func AddMessagePrefix(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

// GenerateKey generates a new ECDSA private key.
func GenerateKey() (*PrivateKey, error) {
	return ecdsa.GenerateKey(s256, rand.Reader)
}

// PublicKeyToAddress returns the Ethereum address for the given ECDSA public key.
func PublicKeyToAddress(publicKey *PublicKey) (addr Address) {
	h := keccak.Hash256(elliptic.Marshal(s256, publicKey.X, publicKey.Y)[1:])
	copy(addr[:], h[12:])
	return
}

// SignHash signs the given hash with the given private key.
func SignHash(privateKey *PrivateKey, hash Hash) (*Signature, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("missing private key")
	}
	privKey, _ := btcec.PrivKeyFromBytes(privateKey.D.Bytes())
	sig, err := btcececdsa.SignCompact(privKey, hash[:], false)
	if err != nil {
		return nil, err
	}
	v := sig[0]
	switch v {
	case 27, 28:
		v -= 27
	}
	copy(sig, sig[1:])
	sig[64] = v
	return &Signature{
		V: new(big.Int).SetBytes(sig[64:]),
		R: new(big.Int).SetBytes(sig[:32]),
		S: new(big.Int).SetBytes(sig[32:64]),
	}, nil
}

// RecoverHash recovers the Ethereum address from the given hash and signature.
func RecoverHash(hash Hash, signature Signature) (*Address, error) {
	if signature.V.BitLen() > 8 {
		return nil, errors.New("invalid signature: V has more than 8 bits")
	}
	if signature.R.BitLen() > 256 {
		return nil, errors.New("invalid signature: R has more than 256 bits")
	}
	if signature.S.BitLen() > 256 {
		return nil, errors.New("invalid signature: S has more than 256 bits")
	}
	v := byte(signature.V.Uint64())
	switch v {
	case 0, 1:
		v += 27
	}
	rb := signature.R.Bytes()
	sb := signature.S.Bytes()
	bin := make([]byte, 65)
	bin[0] = v
	copy(bin[1+(32-len(rb)):], rb)
	copy(bin[33+(32-len(sb)):], sb)
	pub, _, err := btcececdsa.RecoverCompact(bin, hash[:])
	if err != nil {
		return nil, err
	}
	addr := PublicKeyToAddress(pub.ToECDSA())
	return &addr, nil
}

// SignMessage signs the given message with the given private key.
func SignMessage(key *PrivateKey, data []byte) (*Signature, error) {
	if key == nil {
		return nil, fmt.Errorf("missing private key")
	}
	sig, err := SignHash(key, Hash(keccak.Hash256(AddMessagePrefix(data))))
	if err != nil {
		return nil, err
	}
	sig.V = new(big.Int).Add(sig.V, big.NewInt(27))
	return sig, nil
}

// RecoverMessage recovers the Ethereum address from the given message and signature.
func RecoverMessage(data []byte, sig Signature) (*Address, error) {
	sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
	return RecoverHash(Hash(keccak.Hash256(AddMessagePrefix(data))), sig)
}
