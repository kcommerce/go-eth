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

func TestHijackLegacyGasFee(t *testing.T) {
	tt := []struct {
		name     string
		gasLimit *hijackLegacyGasFee
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set gas price",
			gasLimit: &hijackLegacyGasFee{multiplier: 1.0},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionLegacy()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"gasPrice":"0x1000"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
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
			hijacker := transport.NewHijacker(httpMock, tc.gasLimit)

			var result any
			err := hijacker.Call(ctx, &result, tc.method, tc.args...)
			assert.Len(t, tc.request, 0)
			assert.Len(t, tc.response, 0)
			require.NoError(t, err)
		})
	}
}

func TestHijackDynamicGasFee(t *testing.T) {
	tt := []struct {
		name     string
		gasLimit *hijackDynamicGasFee
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set gas price",
			gasLimit: &hijackDynamicGasFee{gasPriceMultiplier: 1.0, priorityFeePerGasMultiplier: 1.0},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionDynamicFee()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_maxPriorityFeePerGas","params":[]}`,
				`{"jsonrpc":"2.0","id":3,"method":"eth_sendTransaction","params":[{"maxFeePerGas":"0x1000","maxPriorityFeePerGas":"0x100"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x100"}`,
				`{"jsonrpc":"2.0","id":3,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
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
			hijacker := transport.NewHijacker(httpMock, tc.gasLimit)

			var result any
			err := hijacker.Call(ctx, &result, tc.method, tc.args...)
			assert.Len(t, tc.request, 0)
			assert.Len(t, tc.response, 0)
			require.NoError(t, err)
		})
	}
}
