package types

import "encoding/json"

type Call interface {
	json.Marshaler
	json.Unmarshaler

	HasCallData
}
