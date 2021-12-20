package idcard

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckBefore(t *testing.T) {
	assert := assert.New(t)

	idnumbers := "12121212"
	err := CheckBefore(idnumbers)
	assert.Error(err)

	idnumbers = "0123456789123456789"
	err = CheckBefore(idnumbers)
	assert.Error(err)

	idnumbers = "012345678912345"
	err = CheckBefore(idnumbers)
	assert.NoError(err)

	idnumbers = "012345678912345678"
	err = CheckBefore(idnumbers)
	assert.NoError(err)
}

func TestClean(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		assert := assert.New(t)
		idnumbers := strconv.FormatInt(time.Now().UnixNano(), 10)
		assert.NotEmpty(idnumbers)
		filtered, isChanged, err := Clean(idnumbers)
		assert.NoError(err)
		assert.Equal(false, isChanged)
		assert.Equal(idnumbers, filtered)
	})

	for _, r := range []rune{'x', 'ｘ', 'Ｘ'} {
		t.Run("cleaned/"+string(r), func(t *testing.T) {
			assert := assert.New(t)
			idnumbers := "0000000010101010" + string(r)
			filtered, isChanged, err := Clean(idnumbers)
			assert.NoError(err)
			assert.Equal(true, isChanged)
			assert.Equal("0000000010101010X", filtered)
		})
	}
}
