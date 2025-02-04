package graceful

import (
	"context"
	"errors"
	"testing"
)

func TestGroupNoTasksErr(t *testing.T) {
	var gr Group

	ctx, cancel := context.WithCancel(context.Background())

	go cancel()

	err := gr.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled: %v", err)
	}
}

func TestGroupNoTasksNoErr(t *testing.T) {
	var gr Group

	ctx, cancel := context.WithCancel(context.Background())

	go cancel()

	err := gr.Run(ctx, IgnoreErrors(context.Canceled))
	if err != nil {
		t.Errorf("should have been ignored: %v", err)
	}
}

func TestGroupNoError(t *testing.T) {
	var gr Group
	//	gr := New()
	//	gr.Signals = nil

	gr.Add(func(ctx context.Context) error {
		return nil
	})

	err := gr.Run(context.Background())
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled: %v", err)
	}
}

func TestGroupNoErrorNoCanceled(t *testing.T) {
	var gr Group
	//	gr := New()
	//	gr.Signals = nil

	gr.Add(func(ctx context.Context) error {
		return nil
	})

	err := gr.Run(context.Background(), IgnoreErrors(context.Canceled))
	if err != nil {
		t.Errorf("should have been ignored: %v", err)
	}
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
	if !errors.Is(err, e) {
		t.Errorf("expected %v, got %v", e, err)
	}
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
	if !errors.Is(err, e) {
		t.Errorf("expected %v, got %v", e, err)
	}
}

func TestBaseOption(t *testing.T) {
	o := IgnoreErrors(context.Canceled)

	t.Logf("option: %v", o)
}
