package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionLegacy_RLP(t *testing.T) {
	tests := []struct {
		tx   *TransactionLegacy
		want []byte
	}{
		{
			tx:   &TransactionLegacy{},
			want: hexutil.MustHexToBytes("c9808080808080808080"),
		},
		{
			tx: &TransactionLegacy{
				EmbedTransactionData: EmbedTransactionData{
					Nonce:     ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				EmbedCallData: EmbedCallData{
					From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
					To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(100000)),
					Input:    []byte{1, 2, 3, 4},
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(1000000000),
				},
			},
			want: hexutil.MustHexToBytes("f87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Example from EIP-155:
		{
			tx: &TransactionLegacy{
				EmbedTransactionData: EmbedTransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(9)),
					Signature: SignatureFromVRSPtr(
						func() *big.Int {
							v, _ := new(big.Int).SetString("37", 10)
							return v
						}(),
						func() *big.Int {
							v, _ := new(big.Int).SetString("18515461264373351373200002665853028612451056578545711640558177340181847433846", 10)
							return v
						}(),
						func() *big.Int {
							v, _ := new(big.Int).SetString("46948507304638947509940763649030358759909902576025900602547168820602576006531", 10)
							return v
						}(),
					),
				},
				EmbedCallData: EmbedCallData{
					To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(21000)),
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(20000000000),
				},
			},
			want: hexutil.MustHexToBytes("f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			// Encode:
			rlp, err := tt.tx.EncodeRLP()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, rlp)

			// Decode:
			tx := NewTransactionLegacy()
			_, err = tx.DecodeRLP(rlp)
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			require.NoError(t1, err)
			equalTX(t1, tx, tt.tx)
		})
	}
}

func TestTransactionLegacy_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		tx   *TransactionLegacy
		want Hash
	}{
		{
			tx:   &TransactionLegacy{},
			want: MustHashFromHex("0x5460be86ce1e4ca0564b5761c6e7070d9f054b671f5404268335000806423d75", PadNone),
		},
		{
			tx: &TransactionLegacy{
				EmbedTransactionData: EmbedTransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				EmbedCallData: EmbedCallData{
					From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
					To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(100000)),
					Input:    []byte{1, 2, 3, 4},
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(1000000000),
				},
			},
			want: MustHashFromHex("0x1efbe489013ac8c0dad2202f68ac12657471df8d80f70e0683ec07b0564a32ca", PadNone),
		},
		// Example from EIP-155:
		{
			tx: &TransactionLegacy{
				EmbedTransactionData: EmbedTransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(9)),
				},
				EmbedCallData: EmbedCallData{
					To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(21000)),
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(20000000000),
				},
			},
			want: MustHashFromHex("0xdaf5a779ae972f972197303d7b574746c7ef83eadac0f2791ad23db92e4c8e53", PadNone),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			sh, err := tt.tx.CalculateSigningHash(keccak256)
			require.NoError(t1, err)
			require.Equal(t1, tt.want, sh)
		})
	}
}
