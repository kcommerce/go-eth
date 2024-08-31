package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionAccessList_RLP(t *testing.T) {
	tests := []struct {
		tx   *TransactionAccessList
		want []byte
	}{
		{
			tx:   &TransactionAccessList{},
			want: hexutil.MustHexToBytes("01cb80808080808080c0808080"),
		},
		{
			tx: &TransactionAccessList{
				EmbedTransactionData: EmbedTransactionData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
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
				EmbedAccessListData: EmbedAccessListData{
					AccessList: []AccessTuple{{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
						},
					}},
				},
			},
			want: hexutil.MustHexToBytes("01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			// Encode:
			rlp, err := tt.tx.EncodeRLP()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, rlp)

			// Decode:
			tx := NewTransactionAccessList()
			_, err = tx.DecodeRLP(rlp)
			tx.From = tt.tx.From
			require.NoError(t1, err)
			equalTX(t1, tx, tt.tx)
		})
	}
}

func TestTransactionAccessList_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		tx   *TransactionAccessList
		want Hash
	}{
		{
			tx:   &TransactionAccessList{},
			want: MustHashFromHex("0xc0157440e7609b2ddee74686831421f05b238ed4c981363e64df8eb1c1ea6afc", PadNone),
		},
		{
			tx: &TransactionAccessList{
				EmbedTransactionData: EmbedTransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				EmbedCallData: EmbedCallData{
					To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(100000)),
					Input:    []byte{1, 2, 3, 4},
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(1000000000),
				},
			},
			want: MustHashFromHex("0x46ba790cdf341de06f08944eecd84721e9ae3c4324098f882597d9817eeba63b", PadNone),
		},
		{
			tx: &TransactionAccessList{
				EmbedTransactionData: EmbedTransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				EmbedCallData: EmbedCallData{
					To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
					Value:    big.NewInt(1000000000000000000),
					GasLimit: ptr(uint64(100000)),
					Input:    []byte{1, 2, 3, 4},
				},
				EmbedLegacyPriceData: EmbedLegacyPriceData{
					GasPrice: big.NewInt(1000000000),
				},
				EmbedAccessListData: EmbedAccessListData{
					AccessList: AccessList{
						AccessTuple{
							Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
							StorageKeys: []Hash{
								MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
								MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
							},
						},
					},
				},
			},
			want: MustHashFromHex("0x71cba0039a020b7a524d7746b79bf6d1f8a521eb1a76715d00116ef1c0f56107", PadNone),
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
