package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func TestHijackAddress(t *testing.T) {
	tt := []struct {
		name     string
		chainID  *hijackAddress
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:    "set address",
			chainID: &hijackAddress{replace: false, address: types.MustAddressFromHex("0x1111111111111111111111111111111111111111")},
			method:  "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetFrom(types.MustAddressFromHex("0x2222222222222222222222222222222222222222"))
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"from": "0x1111111111111111111111111111111111111111"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:    "replace address",
			chainID: &hijackAddress{replace: true, address: types.MustAddressFromHex("0x1111111111111111111111111111111111111111")},
			method:  "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetFrom(types.MustAddressFromHex("0x2222222222222222222222222222222222222222"))
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"from": "0x2222222222222222222222222222222222222222"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			httpMock := newHTTPMock()
			httpMock.Handler = func(req *http.Request) (*http.Response, error) {
				require.NotEmpty(t, tc.request)
				require.NotEmpty(t, tc.response)

				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.JSONEq(t, tc.request[0], string(body), fmt.Sprintf("expected: %s, got: %s", tc.request[0], string(body)))

				res := tc.response[0]
				tc.request = tc.request[1:]
				tc.response = tc.response[1:]
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(res)),
				}, nil
			}
			hijacker := transport.NewHijacker(httpMock, tc.chainID)

			var result any
			err := hijacker.Call(ctx, &result, tc.method, tc.args...)
			assert.Len(t, tc.request, 0)
			assert.Len(t, tc.response, 0)
			require.NoError(t, err)
		})
	}
}
