package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// subscribe creates a subscription to the given method and returns a channel
// that will receive the subscription messages. The messages are unmarshalled
// to the T type. The subscription is unsubscribed and channel closed when the
// context is cancelled.
func subscribe[T any](ctx context.Context, t transport.Transport, method string, params ...any) (chan T, error) {
	st, ok := t.(transport.SubscriptionTransport)
	if !ok {
		return nil, errors.New("transport does not support subscriptions")
	}
	rawCh, subID, err := st.Subscribe(ctx, method, params...)
	if err != nil {
		return nil, err
	}
	msgCh := make(chan T)
	go func() {
		defer close(msgCh)
		defer st.Unsubscribe(ctx, subID)
		for {
			select {
			case <-ctx.Done():
				return
			case raw, ok := <-rawCh:
				if !ok {
					return
				}
				var msg T
				if err := json.Unmarshal(raw, &msg); err != nil {
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, nil
}

// signTransactionResult is the result of an eth_signTransaction request.
// Some backends return only RLP encoded data, others return a JSON object,
// this type can handle both.
type signTransactionResult struct {
	Raw types.Bytes `json:"raw"`
}

func (s *signTransactionResult) UnmarshalJSON(input []byte) error {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		return json.Unmarshal(input, &s.Raw)
	}
	var dec signTransactionResult
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	s.Raw = dec.Raw
	return nil
}
