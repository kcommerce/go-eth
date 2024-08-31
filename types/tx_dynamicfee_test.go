package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionDynamicFee_RLP(t *testing.T) {
	tests := []struct {
		tx   *TransactionDynamicFee
		want []byte
	}{
		{
			tx:   &TransactionDynamicFee{},
			want: hexutil.MustHexToBytes("02cc8080808080808080c0808080"),
		},
		{
			tx: &TransactionDynamicFee{
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
			},
			want: hexutil.MustHexToBytes("02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			// Encode:
			rlp, err := tt.tx.EncodeRLP()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, rlp)

			// Decode:
			tx := NewTransactionDynamicFee()
			_, err = tx.DecodeRLP(rlp)
			tx.From = tt.tx.From
			require.NoError(t1, err)
			equalTX(t1, tx, tt.tx)
		})
	}
}

func TestTransactionDynamicFee_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		tx   *TransactionDynamicFee
		want Hash
	}{
		{
			tx:   &TransactionDynamicFee{},
			want: MustHashFromHex("0x292edeba1be7c90f4dbaed50c44b7f6378633f933202ffe4f547e5a5c2ca3304", PadNone),
		},
		{
			tx: &TransactionDynamicFee{
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
			want: MustHashFromHex("0xc3266152306909bfe339f90fad4f73f958066860300b5a22b98ee6a1d629706c", PadNone),
		},
		{
			tx: &TransactionDynamicFee{
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
			},
			want: MustHashFromHex("0xa66ab756479bfd56f29658a8a199319094e84711e8a2de073ec136ef5179c4c9", PadNone),
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
