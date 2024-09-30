package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

const mockNewFilterRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_newFilter",
	  "params": [
		{
		  "fromBlock": "0x1",
		  "toBlock": "0x2",
		  "address": "0x3333333333333333333333333333333333333333",
		  "topics": ["0x4444444444444444444444444444444444444444444444444444444444444444"]
		}
	  ]
	}
`

const mockNewFilterResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_NewFilter(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockNewFilterRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockNewFilterResponse)),
		}, nil
	}

	from := types.MustBlockNumberFromHex("0x1")
	to := types.MustBlockNumberFromHex("0x2")
	id, err := client.NewFilter(context.Background(), &types.FilterLogsQuery{
		FromBlock: &from,
		ToBlock:   &to,
		Address:   []types.Address{types.MustAddressFromHex("0x3333333333333333333333333333333333333333")},
		Topics: [][]types.Hash{
			{types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone)},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "1", id.String())
}
