package wildcard

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_next(t *testing.T) {
	assert.Equal(t, "b", next("a"))
	assert.Equal(t, "ab", next("aa"))
	assert.Equal(t, "aa", next("z"))
	assert.Equal(t, "ba", next("az"))
	assert.Equal(t, "aaaba", next("aaaaz"))
	assert.Equal(t, "aaa", next("zz"))
}
