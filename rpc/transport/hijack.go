package transport

import (
	"context"
	"encoding/json"
)

const hijackerKey = "goeth_transport_hijackers"

type (
	// Hijacker is used to intercept and modify calls to the underlying
	// Transport. It employs the middleware pattern to facilitate the use of
	// multiple hijackers. The 'next' function should be invoked to proceed
	// with the call chain. The transport provided to the 'next' function is an
	// instance of the underlying Transport. Using it will bypass any
	// subsequent hijackers.
	Hijacker interface {
		// Call returns a CallFunc that intercepts and modifies the Call
		// method. If nil is returned, the Call method is not modified.
		Call() func(next CallFunc) CallFunc

		// Subscribe returns a SubscribeFunc that intercepts and modifies the
		// Subscribe method. If nil is returned, the Subscribe method is not
		// modified.
		Subscribe() func(next SubscribeFunc) SubscribeFunc

		// Unsubscribe returns an UnsubscribeFunc that intercepts and modifies
		// the Unsubscribe method. If nil is returned, the Unsubscribe method
		// is not modified.
		Unsubscribe() func(next UnsubscribeFunc) UnsubscribeFunc
	}

	CallFunc        func(ctx context.Context, t Transport, result any, method string, args ...any) error
	SubscribeFunc   func(ctx context.Context, t SubscriptionTransport, method string, args ...any) (ch chan json.RawMessage, id string, err error)
	UnsubscribeFunc func(ctx context.Context, t SubscriptionTransport, id string) error
)

// Hijack is a wrapper around another Transport that allows for hijacking
// and modifying the behavior of the underlying Transport.
//
// To use Hijack, create a new Hijack instance with the underlying Transport
// and then use the Use method to add any number of hijackers. The hijackers
// will be called in the order they are added.
//
// Hijackers must implement one or more of the Hijacker, SubscribeHijacker,
// and UnsubscribeHijacker interfaces.
type Hijack struct {
	transport Transport
	callFunc  CallFunc
	subFunc   SubscribeFunc
	unsubFunc UnsubscribeFunc
}

func NewHijacker(t Transport, hs ...Hijacker) *Hijack {
	h := &Hijack{
		transport: t,
		callFunc:  defCall,
		subFunc:   defSub,
		unsubFunc: defUnsub,
	}
	h.Use(hs...)
	return h
}

func (h *Hijack) Use(hs ...Hijacker) {
	for _, m := range hs {
		if m == nil {
			continue
		}
		if call := m.Call(); call != nil {
			h.callFunc = call(h.callFunc)
		}
		if sub := m.Subscribe(); sub != nil {
			h.subFunc = sub(h.subFunc)
		}
		if unsub := m.Unsubscribe(); unsub != nil {
			h.unsubFunc = unsub(h.unsubFunc)
		}
	}
}

func (h *Hijack) Call(ctx context.Context, result any, method string, args ...any) error {
	hs := getHijackers(ctx)
	if len(hs) == 0 {
		return h.callFunc(ctx, h.transport, result, method, args...)
	}
	fn := h.callFunc
	for i := len(hs) - 1; i >= 0; i-- {
		if call := hs[i].Call(); call != nil {
			fn = call(fn)
		}
	}
	return fn(ctx, h.transport, result, method, args...)
}

func (h *Hijack) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	if s, ok := h.transport.(SubscriptionTransport); ok {
		return h.subFunc(ctx, s, method, args...)
	}
	return nil, "", ErrNotSubscriptionTransport
}

func (h *Hijack) Unsubscribe(ctx context.Context, id string) error {
	if s, ok := h.transport.(SubscriptionTransport); ok {
		return h.unsubFunc(ctx, s, id)
	}
	return ErrNotSubscriptionTransport
}

func WithHijackers(ctx context.Context, hs ...Hijacker) context.Context {
	return context.WithValue(ctx, hijackerKey, append(getHijackers(ctx), hs...))
}

func getHijackers(ctx context.Context) []Hijacker {
	if hs, ok := ctx.Value(hijackerKey).([]Hijacker); ok {
		return hs
	}
	return nil
}

func defCall(ctx context.Context, t Transport, result any, method string, args ...any) error {
	return t.Call(ctx, result, method, args...)
}

func defSub(ctx context.Context, t SubscriptionTransport, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	return t.Subscribe(ctx, method, args...)
}

func defUnsub(ctx context.Context, t SubscriptionTransport, id string) error {
	return t.Unsubscribe(ctx, id)
}
