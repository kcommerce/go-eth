package types

import "encoding/json"

type CallAccessList struct {
	EmbedCallData
	EmbedLegacyPriceData
	EmbedAccessListData
}

func NewCallAccessList() *CallAccessList {
	return &CallAccessList{}
}

func (c *CallAccessList) Copy() *CallAccessList {
	if c == nil {
		return nil
	}
	return &CallAccessList{
		EmbedCallData:        *c.EmbedCallData.Copy(),
		EmbedLegacyPriceData: *c.EmbedLegacyPriceData.Copy(),
		EmbedAccessListData:  *c.EmbedAccessListData.Copy(),
	}
}

func (c *CallAccessList) MarshalJSON() ([]byte, error) {
	call := &jsonCallAccessList{
		From:       c.From,
		To:         c.To,
		Data:       c.Input,
		AccessList: c.AccessList,
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

func (c *CallAccessList) UnmarshalJSON(data []byte) error {
	call := &jsonCallAccessList{}
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
	c.AccessList = call.AccessList
	return nil
}

type jsonCallAccessList struct {
	From       *Address   `json:"from,omitempty"`
	To         *Address   `json:"to,omitempty"`
	GasLimit   *Number    `json:"gas,omitempty"`
	GasPrice   *Number    `json:"gasPrice,omitempty"`
	Value      *Number    `json:"value,omitempty"`
	Data       Bytes      `json:"data,omitempty"`
	AccessList AccessList `json:"accessList,omitempty"`
}
