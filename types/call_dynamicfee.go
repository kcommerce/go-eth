package types

import "encoding/json"

type CallDynamicFee struct {
	EmbedCallData
	EmbedDynamicFeeData
	EmbedAccessListData
}

func NewCallDynamicFee() *CallDynamicFee {
	return &CallDynamicFee{}
}

func (c *CallDynamicFee) Copy() *CallDynamicFee {
	if c == nil {
		return nil
	}
	return &CallDynamicFee{
		EmbedCallData:       *c.EmbedCallData.Copy(),
		EmbedDynamicFeeData: *c.EmbedDynamicFeeData.Copy(),
		EmbedAccessListData: *c.EmbedAccessListData.Copy(),
	}
}

func (c *CallDynamicFee) MarshalJSON() ([]byte, error) {
	call := &jsonCallDynamicFee{
		From:       c.From,
		To:         c.To,
		Data:       c.Input,
		AccessList: c.AccessList,
	}
	if c.GasLimit != nil {
		call.GasLimit = NumberFromUint64Ptr(*c.GasLimit)
	}
	if c.MaxFeePerGas != nil {
		call.MaxFeePerGas = NumberFromBigIntPtr(c.MaxFeePerGas)
	}
	if c.MaxPriorityFeePerGas != nil {
		call.MaxPriorityFeePerGas = NumberFromBigIntPtr(c.MaxPriorityFeePerGas)
	}
	if c.Value != nil {
		value := NumberFromBigInt(c.Value)
		call.Value = &value
	}
	return json.Marshal(call)
}

func (c *CallDynamicFee) UnmarshalJSON(data []byte) error {
	call := &jsonCallDynamicFee{}
	if err := json.Unmarshal(data, call); err != nil {
		return err
	}
	c.From = call.From
	c.To = call.To
	if call.GasLimit != nil {
		gas := call.GasLimit.Big().Uint64()
		c.GasLimit = &gas
	}
	if call.MaxFeePerGas != nil {
		c.MaxFeePerGas = call.MaxFeePerGas.Big()
	}
	if call.MaxPriorityFeePerGas != nil {
		c.MaxPriorityFeePerGas = call.MaxPriorityFeePerGas.Big()
	}
	if call.Value != nil {
		c.Value = call.Value.Big()
	}
	c.Input = call.Data
	c.AccessList = call.AccessList
	return nil
}

type jsonCallDynamicFee struct {
	From                 *Address   `json:"from,omitempty"`
	To                   *Address   `json:"to,omitempty"`
	GasLimit             *Number    `json:"gas,omitempty"`
	MaxFeePerGas         *Number    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *Number    `json:"maxPriorityFeePerGas,omitempty"`
	Value                *Number    `json:"value,omitempty"`
	Data                 Bytes      `json:"data,omitempty"`
	AccessList           AccessList `json:"accessList,omitempty"`
}
