package dashboard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseActions(t *testing.T) {
	raw := json.RawMessage(`[{"action_type": "purchase", "value": "10"}]`)
	actions := parseActions(raw)

	assert.Len(t, actions, 1)
	assert.Equal(t, "purchase", actions[0].ActionType)
	assert.Equal(t, "10", actions[0].Value)

	assert.Nil(t, parseActions(nil))
}

func TestFindAction(t *testing.T) {
	actions := []metaAction{
		{ActionType: "purchase", Value: "10"},
		{ActionType: "link_click", Value: "20"},
	}

	assert.Equal(t, "10", findAction(actions, "purchase"))
	assert.Equal(t, "20", findAction(actions, "link_click"))
	assert.Equal(t, "0", findAction(actions, "unknown"))
}

func TestFormatFloat(t *testing.T) {
	assert.Equal(t, "0", formatFloat(0))
	assert.Equal(t, "10.5", formatFloat(10.5))
}

func TestFormatNullFloat(t *testing.T) {
	assert.Equal(t, "0", formatNullFloat(nil))
	val := 10.5
	assert.Equal(t, "10.5", formatNullFloat(&val))
}

func TestFormatNullInt(t *testing.T) {
	assert.Equal(t, "0", formatNullInt(nil))
	val := int64(100)
	assert.Equal(t, "100", formatNullInt(&val))
}

func TestResolveBudget(t *testing.T) {
	assert.Equal(t, "100", resolveBudget(100, 0))
	assert.Equal(t, "500", resolveBudget(0, 500))
	assert.Equal(t, "100", resolveBudget(100, 500))
	assert.Equal(t, "0", resolveBudget(0, 0))
}
