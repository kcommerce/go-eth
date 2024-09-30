package types

import "encoding/json"

type CallLegacy struct {
	EmbedCallData
	EmbedLegacyPriceData
}

func NewCallLegacy() *CallLegacy {
	return &CallLegacy{}
}

func (c *CallLegacy) Copy() *CallLegacy {
	if c == nil {
		return nil
	}
	return &CallLegacy{
		EmbedCallData:        *c.EmbedCallData.Copy(),
		EmbedLegacyPriceData: *c.EmbedLegacyPriceData.Copy(),
	}
}

func (c CallLegacy) MarshalJSON() ([]byte, error) {
	call := &jsonCallLegacy{
		From: c.From,
		To:   c.To,
		Data: c.Input,
	}
	if c.GasLimit != nil {
		call.GasLimit = NumberFromUint64Ptr(*c.GasLimit)
	}
	if c.GasPrice != nil {
		call.GasPrice = NumberFromBigIntPtr(c.GasPrice)
	}
	if c.Value != nil {
		value := NumberFromBigInt(c.Value)
		call.Value = &value
	}
	return json.Marshal(call)
}

func (c *CallLegacy) UnmarshalJSON(data []byte) error {
	call := &jsonCallLegacy{}
	if err := json.Unmarshal(data, call); err != nil {
		return err
	}
	c.From = call.From
	c.To = call.To
	if call.GasLimit != nil {
		gas := call.GasLimit.Big().Uint64()
		c.GasLimit = &gas
	}
	if call.GasPrice != nil {
		c.GasPrice = call.GasPrice.Big()
	}
	if call.Value != nil {
		c.Value = call.Value.Big()
	}
	c.Input = call.Data
	return nil
}

type jsonCallLegacy struct {
	From     *Address `json:"from,omitempty"`
	To       *Address `json:"to,omitempty"`
	GasLimit *Number  `json:"gas,omitempty"`
	GasPrice *Number  `json:"gasPrice,omitempty"`
	Value    *Number  `json:"value,omitempty"`
	Data     Bytes    `json:"data,omitempty"`
}
