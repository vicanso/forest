package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert := assert.New(t)
		s := &struct {
			Method string `json:"method" valid:"xMethods"`
			Type   int    `json:"type" valid:"xIntIn(1|3|5)"`
			Size   int    `json:"size" valid:"xIntRange(1|2)"`
		}{}
		err := Do(s, []byte(`{
			"method": "GET",
			"type": 3,
			"size": 10
		}`))
		assert.NotNil(err)
		err = Do(s, map[string]interface{}{
			"method": "GET",
			"type":   3,
			"size":   1,
		})
		assert.Nil(err)
	})
}
