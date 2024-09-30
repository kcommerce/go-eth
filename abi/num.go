package abi

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

var (
	// MaxUint contains the maximum unsigned integer for each bit size.
	MaxUint = map[int]*big.Int{}

	// MaxInt contains the maximum signed integer for each bit size.
	MaxInt = map[int]*big.Int{}

	// MinInt contains the minimum signed integer for each bit size.
	MinInt = map[int]*big.Int{}
)

// intX represents a signed integer of bit size between 8 and 256.
type intX struct {
	size int
	val  *big.Int
}

// newIntX creates a new intX value.
func newIntX(bitSize int) *intX {
	if bitSize < 8 || bitSize > 256 || bitSize%8 != 0 {
		panic("abi: invalid bit size for intX")
	}
	return &intX{
		size: bitSize,
		val:  new(big.Int),
	}
}

// BitSize returns the bit size of the integer.
func (i *intX) BitSize() int {
	return i.size
}

// BitLen returns the number of bits required to represent the integer.
func (i *intX) BitLen() int {
	return signedBitLen(i.val)
}

// IsInt returns true if the integer can be represented as an int.
func (i *intX) IsInt() bool {
	if !i.val.IsInt64() {
		return false
	}
	x := i.val.Int64()
	if x > math.MaxInt {
		return false
	}
	if x < math.MinInt {
		return false
	}
	return true
}

// Int returns the int64 representation of the integer.
func (i *intX) Int() (int, error) {
	if !i.val.IsInt64() {
		return 0, fmt.Errorf("abi: int overflow")
	}
	x := i.val.Int64()
	if x > math.MaxInt {
		return 0, fmt.Errorf("abi: int overflow")
	}
	if x < math.MinInt {
		return 0, fmt.Errorf("abi: int overflow")
	}
	return int(i.val.Int64()), nil
}

// Int64 returns the int64 representation of the integer.
func (i *intX) Int64() (int64, error) {
	if !i.val.IsInt64() {
		return 0, fmt.Errorf("abi: int64 overflow")
	}
	return i.val.Int64(), nil
}

// BigInt returns the *big.Int representation of the integer.
func (i *intX) BigInt() *big.Int {
	return i.val
}

// Bytes returns the value of the integer as a big-endian byte slice.
// The byte slice is zero-padded to the size of the integer. Negative
// values are two's complement encoded.
func (i *intX) Bytes() []byte {
	r := make([]byte, i.size/8)
	x := new(big.Int).Set(i.val).And(i.val, MaxUint[i.size])
	padLeft(r, x.Bytes())
	return r
}

// FillBytes fills the byte slice r with the value of the integer as a
// big-endian byte slice. If the byte slice is smaller than the integer,
// an error is returned. The byte slice cannot be larger than 256 bits.
// Negative values are two's complement encoded.
func (i *intX) FillBytes(r []byte) error {
	bitLen := len(r) * 8
	if bitLen < i.size {
		return fmt.Errorf("abi: cannot fill %d-bit integer to %d bytes", i.size, len(r))
	}
	if bitLen > 256 {
		return fmt.Errorf("abi: cannot fill byte slices larger than 256 bits")
	}
	x := new(big.Int).Set(i.val).And(i.val, MaxUint[bitLen])
	x.FillBytes(r)
	return nil
}

// SetInt sets the value of the integer to x.
func (i *intX) SetInt(x int) error {
	if bits.Len(uint(x)) > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit int", bits.Len(uint(x)), i.size)
	}
	i.val.SetInt64(int64(x))
	return nil
}

// SetInt64 sets the value of the integer to x.
func (i *intX) SetInt64(x int64) error {
	if bits.Len64(uint64(x)) > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit int64", bits.Len64(uint64(x)), i.size)
	}
	i.val.SetInt64(x)
	return nil
}

// SetBigInt sets the value of the integer to x. If x is larger than the
// integer's bit size, an error is returned.
func (i *intX) SetBigInt(x *big.Int) error {
	if x == nil || x.Sign() == 0 {
		i.val = big.NewInt(0)
		return nil
	}
	if signedBitLen(x) > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit signed int", signedBitLen(x), i.size)
	}
	i.val.Set(x)
	return nil
}

// SetBytes sets the value of the integer to x. If x is larger than the
// integer's bit size, an error is returned.
func (i *intX) SetBytes(b []byte) error {
	x := new(big.Int).SetBytes(b)
	if x.Cmp(MaxInt[i.size]) > 0 {
		// If the number is negative, we need to set it from the two's complement
		// representation.
		x.Add(MaxUint[i.size], new(big.Int).Neg(x))
		x.Add(x, big.NewInt(1))
		x.Neg(x)
	}
	return i.SetBigInt(x)
}

// uintX represents a unsigned integer of bit size between 8 and 256.
type uintX struct {
	size int
	val  *big.Int
}

// newUintX creates a new uintX value.
func newUintX(bitSize int) *uintX {
	if bitSize < 8 || bitSize > 256 || bitSize%8 != 0 {
		panic("abi: invalid bit size for intX")
	}
	return &uintX{
		size: bitSize,
		val:  new(big.Int),
	}
}

// Uint returns the uint representation of the integer.
func (i *uintX) Uint() (int, error) {
	if !i.val.IsUint64() {
		return 0, fmt.Errorf("abi: uint overflow")
	}
	x := i.val.Uint64()
	if x > math.MaxUint {
		return 0, fmt.Errorf("abi: uint overflow")
	}
	return int(i.val.Uint64()), nil
}

// Uint64 returns the uint64 representation of the integer.
func (i *uintX) Uint64() (uint64, error) {
	if !i.val.IsUint64() {
		return 0, fmt.Errorf("abi: int64 overflow")
	}
	return i.val.Uint64(), nil
}

// BigInt returns the *big.Int representation of the integer.
func (i *uintX) BigInt() *big.Int {
	return i.val
}

// Bytes returns the value of the integer as a big-endian byte slice.
// The byte slice is zero-padded to the size of the integer. Negative
// values are two's complement encoded.
func (i *uintX) Bytes() []byte {
	r := make([]byte, i.size/8)
	padLeft(r, i.val.Bytes())
	return r
}

// FillBytes fills the byte slice r with the value of the integer as a
// big-endian byte slice. If the byte slice is smaller than the integer,
// an error is returned. The byte slice cannot be larger than 256 bits.
func (i *uintX) FillBytes(r []byte) error {
	bitLen := len(r) * 8
	if bitLen < i.size {
		return fmt.Errorf("abi: cannot fill %d-bit integer to %d bytes", i.size, len(r))
	}
	if bitLen > 256 {
		return fmt.Errorf("abi: cannot fill byte slices larger than 256 bits")
	}
	i.val.FillBytes(r)
	return nil
}

// SetUint sets the value of the integer to x.
func (i *uintX) SetUint(x uint) error {
	if bits.Len(x) > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit int", bits.Len(x), i.size)
	}
	i.val.SetUint64(uint64(x))
	return nil
}

// SetUint64 sets the value of the integer to x.
func (i *uintX) SetUint64(x uint64) error {
	if bits.Len64(x) > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit int64", bits.Len64(x), i.size)
	}
	i.val.SetUint64(x)
	return nil
}

// SetBigInt sets the value of the integer to x. If x is larger than the
// integer's bit size, an error is returned.
func (i *uintX) SetBigInt(x *big.Int) error {
	if x == nil || x.Sign() == 0 {
		i.val = big.NewInt(0)
		return nil
	}
	if x.BitLen() > i.size {
		return fmt.Errorf("abi: cannot set %d-bit integer to %d-bit signed int", signedBitLen(x), i.size)
	}
	i.val.Set(x)
	return nil
}

// SetBytes sets the value of the integer to x. If x is larger than the
// integer's bit size, an error is returned.
func (i *uintX) SetBytes(b []byte) error {
	return i.SetBigInt(new(big.Int).SetBytes(b))
}

func padLeft(dst []byte, src []byte) {
	copy(dst[len(dst)-len(src):], src)
}

// signedBitLen returns the number of bits required to represent x in two's
// complement representation.
func signedBitLen(x *big.Int) int {
	if x == nil || x.Sign() == 0 {
		return 0
	}
	bitLen := x.BitLen()
	if x.Sign() < 0 && x.TrailingZeroBits() == uint(bitLen-1) {
		// If the binary representation of the number is equal to x^2, then the
		// bit length for the negative number encoded in two's complement is
		// one bit shorter.
		return bitLen
	}
	return bitLen + 1
}

func canSetInt(x int64, bitLen int) bool {
	if bitLen >= 64 {
		return true
	}
	if bitLen <= 0 {
		return false
	}
	if x < 0 {
		return x >= -1<<(bitLen-1)
	}
	return x < (1 << uint(bitLen-1))
}

func canSetUint(x uint64, bitLen int) bool {
	if bitLen >= 64 {
		return true
	}
	if bitLen <= 0 {
		return false
	}
	return x < (1 << uint(bitLen))
}

func init() {
	pOne := big.NewInt(1)
	mOne := big.NewInt(-1)
	for i := 8; i <= 256; i += 8 {
		MaxUint[i] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(i)), pOne)
		MaxInt[i] = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(i-1)), pOne)
		MinInt[i] = new(big.Int).Lsh(mOne, uint(i-1))
	}
}
