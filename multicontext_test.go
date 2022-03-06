package graceful

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	done := ctx.Done()

	assert.Nil(t, done) // must not be nil, it's implementation detail
}
