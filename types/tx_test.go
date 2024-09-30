package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func equalTX(t *testing.T, actual, expected Transaction) {
	assert.Equal(t, derefType(reflect.TypeOf(actual)), derefType(reflect.TypeOf(expected)))
	assert.Equal(t, actual.TransactionData(), expected.TransactionData())
	if _, ok := expected.(HasCallData); ok {
		assert.Equal(t, actual.(HasCallData).CallData(), actual.(HasCallData).CallData())
	}
	if _, ok := actual.(HasLegacyPrice); ok {
		assert.Equal(t, expected.(HasLegacyPrice).LegacyPriceData(), actual.(HasLegacyPrice).LegacyPriceData())
	}
	if _, ok := actual.(HasAccessListData); ok {
		assert.Equal(t, expected.(HasAccessListData).AccessListData(), actual.(HasAccessListData).AccessListData())
	}
	if _, ok := actual.(HasDynamicFeeData); ok {
		assert.Equal(t, expected.(HasDynamicFeeData).DynamicFeeData(), actual.(HasDynamicFeeData).DynamicFeeData())
	}
	if _, ok := actual.(HasBlobData); ok {
		assert.Equal(t, expected.(HasBlobData).BlobData(), actual.(HasBlobData).BlobData())
	}
}

func derefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	return t
}
