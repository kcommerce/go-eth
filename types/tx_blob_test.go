package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionBlob_RLP(t *testing.T) {
	tests := []struct {
		tx   *TransactionBlob
		want []byte
	}{
		{
			tx:   &TransactionBlob{},
			want: hexutil.MustHexToBytes("03ce8080808080808080c080c0808080"),
		},
		{
			tx: &TransactionBlob{
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
				EmbedDynamicFeeData: EmbedDynamicFeeData{
					MaxPriorityFeePerGas: big.NewInt(1000000000),
					MaxFeePerGas:         big.NewInt(2000000000),
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
				EmbedBlobData: EmbedBlobData{
					MaxFeePerBlobGas: big.NewInt(3000000000),
					Blobs: []Blob{
						{
							Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
						},
					},
				},
			},
			want: hexutil.MustHexToBytes("03f8fa0101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a0555555555555555555555555555555555555555555555555555555555555555584b2d05e00e1a066666666666666666666666666666666666666666666666666666666666666666fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			// Encode:
			rlp, err := tt.tx.EncodeRLP()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, rlp)

			// Decode:
			tx := NewTransactionBlob()
			_, err = tx.DecodeRLP(rlp)
			tx.From = tt.tx.From
			require.NoError(t1, err)
			equalTX(t1, tx, tt.tx)
		})
	}
}

func TestTransactionBlob_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		tx   *TransactionBlob
		want Hash
	}{
		{
			tx:   &TransactionBlob{},
			want: MustHashFromHex("0x846c9b47f161837f5068b0ffb0c1a98785302f89d613338ccfa9a1c72c9f951d", PadNone),
		},
		{
			tx: &TransactionBlob{
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
				EmbedDynamicFeeData: EmbedDynamicFeeData{
					MaxPriorityFeePerGas: big.NewInt(1000000000),
					MaxFeePerGas:         big.NewInt(2000000000),
				},
			},
			want: MustHashFromHex("0x0604b49731147cf745c666f1a67bf1b5e9fbee127085b3d4c4958191590e8bce", PadNone),
		},
		{
			tx: &TransactionBlob{
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
				EmbedDynamicFeeData: EmbedDynamicFeeData{
					MaxPriorityFeePerGas: big.NewInt(1000000000),
					MaxFeePerGas:         big.NewInt(2000000000),
				},
				EmbedBlobData: EmbedBlobData{
					MaxFeePerBlobGas: big.NewInt(3000000000),
					Blobs: []Blob{
						{
							Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
						},
					},
				},
			},
			want: MustHashFromHex("0x3faa63efab3e460606c31cd9a2e8791d87e91954137571eddb3b4b0abc69e2cd", PadNone),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			sh, err := tt.tx.CalculateSigningHash()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, sh)
		})
	}
}
