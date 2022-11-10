package nuclei

import (
	"context"
	"testing"

	"github.com/projectdiscovery/nuclei/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestWrapper_RunEnumeration(t *testing.T) {
	options := &types.Options{}
	wrapper, _ := NewWrapper(options)
	results, err := wrapper.RunEnumeration(context.Background(), []string{"http://218.77.53.28:8087"})
	assert.Nil(t, err)
	assert.NotNil(t, results)
}
