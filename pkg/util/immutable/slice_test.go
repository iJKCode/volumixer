package immutable

import (
	"github.com/stretchr/testify/assert"
	"maps"
	"slices"
	"testing"
)

func TestSlice(t *testing.T) {
	bad := []int{1, 2, 3, 4}
	good := []int{29, 33, 14, 93, 9, 97, 6, 69, 43, 66}
	slice := SliceFrom(good)

	for idx, val := range good {
		res, ok := slice.Get(idx)
		assert.True(t, ok)
		assert.Equal(t, val, res)
	}

	{
		_, ok := slice.Get(-1)
		assert.False(t, ok)
		_, ok = slice.Get(len(good))
		assert.False(t, ok)
	}

	assert.Equal(t, len(good), slice.Len())
	assert.NotEqual(t, len(bad), slice.Len())

	assert.Equal(t, good, slice.Slice())
	assert.True(t, SliceEqual(slice, SliceFrom(good)))
	assert.False(t, SliceEqual(slice, SliceFrom(bad)))

	{
		clone := slices.Collect(slice.Values())
		assert.Equal(t, good, clone, "iterated values should be equal")
	}
	{
		val := maps.Collect(slice.Items())
		exp := maps.Collect(slices.All(good))
		assert.Equal(t, exp, val, "enumerated values should be equal")
	}
}
