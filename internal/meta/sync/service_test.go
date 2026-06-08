package sync

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBatchCode(t *testing.T) {
	now := time.Now()
	code := generateBatchCode(now)
	assert.True(t, strings.HasPrefix(code, "META-SYNC-"))
}

func TestCalculateDurationMs(t *testing.T) {
	now := time.Now()
	future := now.Add(100 * time.Millisecond)

	duration := calculateDurationMs(&now, &future)
	assert.GreaterOrEqual(t, duration, uint64(100))

	assert.Equal(t, uint64(0), calculateDurationMs(nil, &future))
	assert.Equal(t, uint64(0), calculateDurationMs(&now, nil))
}

func TestNullableString(t *testing.T) {
	str := "test"
	res := nullableString(str)
	assert.NotNil(t, res)
	assert.Equal(t, "test", *res)

	emptyStr := ""
	res2 := nullableString(emptyStr)
	assert.Nil(t, res2)
}
