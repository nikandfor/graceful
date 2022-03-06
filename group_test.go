package graceful

import (
	"context"
	"testing"

	"github.com/nikandfor/assert"
	"github.com/nikandfor/assert/is"
	"github.com/nikandfor/errors"
)

func TestGroupNoError(t *testing.T) {
	gr := New()
	gr.Signals = nil

	gr.Add(context.Background(), "a", func(ctx context.Context) error {
		return nil
	})

	err := gr.Run(context.Background())
	assert.NoError(t, err)
}

func TestGroupSomeError(t *testing.T) {
	gr := New()
	gr.Signals = nil

	e := errors.New("some")

	gr.Add(context.Background(), "a", func(ctx context.Context) error {
		return e
	})

	err := gr.Run(context.Background())
	assert.ErrorIs(t, err, e)
}

func TestStopGlobal(t *testing.T) {
	gr := New()
	gr.Signals = nil

	e := errors.New("some")

	gr.Add(context.Background(), "a", func(ctx context.Context) error {
		<-ctx.Done()
		return e
	})

	gr.Add(context.Background(), "b", func(ctx context.Context) error {
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

func TestStopOne(t *testing.T) {
	gr := New()
	gr.Signals = nil

	ea := errors.New("some a")
	eb := errors.New("some b")

	ctx, cancel := context.WithCancel(context.Background())

	go cancel()

	gr.Add(ctx, "a", func(ctx context.Context) error {
		<-ctx.Done()
		return ea
	})

	gr.Add(context.Background(), "b", func(ctx context.Context) error {
		<-ctx.Done()
		return eb
	})

	err := gr.Run(context.Background())
	assert.ErrorIs(t, err, ea)
}
