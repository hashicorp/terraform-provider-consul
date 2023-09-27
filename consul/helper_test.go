package consul

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatKeysFunc(t *testing.T) {
	dummySetSchema := new(schema.Set)
	dummySetSchema.F = func(i interface{}) int {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("%s-", i))
		return hashcode.String(buf.String())
	}
	dummySetSchema.Add(map[string]interface{}{"test_set_key": "test_value"})
	dataSet := []map[string]interface{}{
		{
			"input": map[string]interface{}{
				"test_key":     "value1",
				"test_tls_key": "value2",
				"ttl":          "value3",
			},
			"expected": map[string]interface{}{
				"TestKey":    "value1",
				"TestTLSKey": "value2",
				"TTL":        "value3",
			},
		},
		{
			"input":    dummySetSchema,
			"expected": map[string]interface{}{"TestSetKey": "test_value"},
		},
	}
	for _, testCase := range dataSet {
		res, err := formatKeys(testCase["input"], formatKey)
		require.NoError(t, err)
		require.Equal(t, testCase["expected"], res)
	}
}
