package types

import "math/big"

// Embedded types are used to embed common fields and methods into call and
// transaction types.

type HasCallData interface {
	CallData() *EmbedCallData
	SetCallData(data *EmbedCallData)
}

type HasTransactionData interface {
	TransactionData() *EmbedTransactionData
	SetTransactionData(data *EmbedTransactionData)
}

type HasLegacyPrice interface {
	LegacyPriceData() *EmbedLegacyPriceData
	SetLegacyPriceData(data *EmbedLegacyPriceData)
}

type HasAccessListData interface {
	AccessListData() *EmbedAccessListData
	SetAccessListData(data *EmbedAccessListData)
}

type HasDynamicFeeData interface {
	DynamicFeeData() *EmbedDynamicFeeData
	SetDynamicFeeData(data *EmbedDynamicFeeData)
}

// EmbedCallData is set of common fields for calls and transactions.
type EmbedCallData struct {
	From     *Address // From is the sender address.
	To       *Address // To is the recipient address. nil means contract creation.
	GasLimit *uint64  // GasLimit is the gas limit, if 0, there is no limit.
	Value    *big.Int // Value is the amount of wei to send.
	Input    []byte   // Input is the input data.
}

func (c *EmbedCallData) CallData() *EmbedCallData {
	return c
}

func (c *EmbedCallData) SetCallData(data *EmbedCallData) {
	*c = *data
}

func (c *EmbedCallData) SetFrom(from Address) {
	c.From = &from
}

func (c *EmbedCallData) SetTo(to Address) {
	c.To = &to
}

func (c *EmbedCallData) SetGasLimit(gasLimit uint64) {
	c.GasLimit = &gasLimit
}

func (c *EmbedCallData) SetValue(value *big.Int) {
	c.Value = value
}

func (c *EmbedCallData) SetInput(input []byte) {
	c.Input = input
}

func (c *EmbedCallData) Copy() *EmbedCallData {
	if c == nil {
		return nil
	}
	return &EmbedCallData{
		From:     copyPtr(c.From),
		To:       copyPtr(c.To),
		GasLimit: copyPtr(c.GasLimit),
		Value:    copyBigInt(c.Value),
		Input:    copyBytes(c.Input),
	}
}

// EmbedTransactionData is a set of common fields for transactions.
type EmbedTransactionData struct {
	ChainID   *uint64
	Nonce     *uint64
	Signature *Signature
}

func (c *EmbedTransactionData) TransactionData() *EmbedTransactionData {
	return c
}

func (c *EmbedTransactionData) SetTransactionData(data *EmbedTransactionData) {
	*c = *data
}

func (c *EmbedTransactionData) SetChainID(chainID uint64) {
	c.ChainID = &chainID
}

func (c *EmbedTransactionData) SetNonce(nonce uint64) {
	c.Nonce = &nonce
}

func (c *EmbedTransactionData) SetSignature(signature Signature) {
	c.Signature = &signature
}

func (c *EmbedTransactionData) Copy() *EmbedTransactionData {
	return &EmbedTransactionData{
		ChainID:   copyPtr(c.ChainID),
		Nonce:     copyPtr(c.Nonce),
		Signature: c.Signature.Copy(),
	}
}

type EmbedLegacyPriceData struct {
	GasPrice *big.Int
}

func (c *EmbedLegacyPriceData) LegacyPriceData() *EmbedLegacyPriceData {
	return c
}

func (c *EmbedLegacyPriceData) SetLegacyPriceData(data *EmbedLegacyPriceData) {
	*c = *data
}

func (c *EmbedLegacyPriceData) SetGasPrice(gasPrice *big.Int) {
	c.GasPrice = gasPrice
}

func (c *EmbedLegacyPriceData) Copy() *EmbedLegacyPriceData {
	if c == nil {
		return nil
	}
	return &EmbedLegacyPriceData{
		GasPrice: copyBigInt(c.GasPrice),
	}
}

type EmbedAccessListData struct {
	AccessList AccessList
}

func (c *EmbedAccessListData) AccessListData() *EmbedAccessListData {
	return c
}

func (c *EmbedAccessListData) SetAccessListData(data *EmbedAccessListData) {
	*c = *data
}

func (c *EmbedAccessListData) SetAccessList(accessList AccessList) {
	c.AccessList = accessList
}

func (c *EmbedAccessListData) Copy() *EmbedAccessListData {
	if c == nil {
		return nil
	}
	return &EmbedAccessListData{
		AccessList: c.AccessList.Copy(),
	}
}

type EmbedDynamicFeeData struct {
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

func (c *EmbedDynamicFeeData) DynamicFeeData() *EmbedDynamicFeeData {
	return c
}

func (c *EmbedDynamicFeeData) SetDynamicFeeData(data *EmbedDynamicFeeData) {
	*c = *data
}

func (c *EmbedDynamicFeeData) SetMaxFeePerGas(maxFeePerGas *big.Int) {
	c.MaxFeePerGas = maxFeePerGas
}

func (c *EmbedDynamicFeeData) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) {
	c.MaxPriorityFeePerGas = maxPriorityFeePerGas
}

func (c *EmbedDynamicFeeData) Copy() *EmbedDynamicFeeData {
	if c == nil {
		return nil
	}
	return &EmbedDynamicFeeData{
		MaxFeePerGas:         copyBigInt(c.MaxFeePerGas),
		MaxPriorityFeePerGas: copyBigInt(c.MaxPriorityFeePerGas),
	}
}
