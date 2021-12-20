package idcard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		assert := assert.New(t)
		for _, idnumbers := range []string{
			"012345678901234567",
			"012345678901234",
			"01234567890123456X",
			"01234567890123X",
		} {
			assert.NotEmpty(idnumbers)
			filtered, err := Clean(idnumbers)
			assert.NoError(err)
			assert.Equal(idnumbers, filtered)
		}
	})

	t.Run("clean spaces", func(t *testing.T) {
		assert := assert.New(t)
		idnumbers := "012345678901234567"
		cleaned, err := Clean("012345678901 2345 67")
		assert.NoError(err)
		assert.Equal(idnumbers, cleaned)
	})

	for _, r := range []rune{'x', 'ｘ', 'Ｘ'} {
		t.Run("cleaned/"+string(r), func(t *testing.T) {
			assert := assert.New(t)
			idnumbers := "01234567890123456" + string(r)
			filtered, err := Clean(idnumbers)
			assert.NoError(err)
			assert.Equal("01234567890123456X", filtered)
		})
	}
}
