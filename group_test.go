package graceful

import (
	"context"
	"testing"

	"github.com/nikandfor/assert"
	"github.com/nikandfor/assert/is"
	"github.com/nikandfor/errors"
)

func TestGroupNoTasks(t *testing.T) {
	var gr Group

	ctx, cancel := context.WithCancel(context.Background())

	go cancel()

	err := gr.Run(ctx)
	assert.NoError(t, err)
}

func TestGroupNoError(t *testing.T) {
	var gr Group
	//	gr := New()
	//	gr.Signals = nil

	gr.Add(func(ctx context.Context) error {
		return nil
	})

	err := gr.Run(context.Background())
	assert.NoError(t, err)
}

func TestGroupSomeError(t *testing.T) {
	var gr Group
	//	gr := New()
	//	gr.Signals = nil

	e := errors.New("some")

	gr.Add(func(ctx context.Context) error {
		return e
	})

	err := gr.Run(context.Background())
	assert.ErrorIs(t, err, e)
}

func TestStopGlobal(t *testing.T) {
	var gr Group
	//	gr := New()
	//	gr.Signals = nil

	e := errors.New("some")

	gr.Add(func(ctx context.Context) error {
		<-ctx.Done()
		return e
	})

	gr.Add(func(ctx context.Context) error {
		<-ctx.Done()
		return e
	})

	ctx, cancel := context.WithCancel(context.Background())

	go cancel()

	err := gr.Run(ctx)
	assert.Any(t, []assert.Checker{
		is.ErrorIs(err, e),
		is.ErrorIs(err, context.Canceled),
	})
}
