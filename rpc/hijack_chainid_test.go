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

const mockHijackChainIDCallValidResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x01"
	}
`

const mockHijackChainIDSendTransactionValidResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}
`

func TestHijackChainID(t *testing.T) {
	tt := []struct {
		name     string
		chainID  *hijackChainID
		method   string
		args     []any
		request  []string
		response []*http.Response
	}{
		{
			name:    "set chainID",
			chainID: &hijackChainID{},
			method:  "eth_sendTransaction",
			args:    []any{types.NewTransactionAccessList()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId": "0x1"}]}`,
			},
			response: []*http.Response{
				{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockHijackChainIDCallValidResponse)),
				},
				{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockHijackChainIDSendTransactionValidResponse)),
				},
			},
		},
		{
			name:    "do not replace chainID",
			chainID: &hijackChainID{replace: false},
			method:  "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetChainID(2)
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"chainId": "0x2"}]}`,
			},
			response: []*http.Response{
				{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockHijackChainIDSendTransactionValidResponse)),
				},
			},
		},
		{
			name:    "replace chainID",
			chainID: &hijackChainID{replace: true},
			method:  "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetChainID(2)
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId": "0x1"}]}`,
			},
			response: []*http.Response{
				{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockHijackChainIDCallValidResponse)),
				},
				{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockHijackChainIDSendTransactionValidResponse)),
				},
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
				return res, nil
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
