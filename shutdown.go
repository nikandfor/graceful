package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

type (
	settings struct {
		Signals []os.Signal

		KillErr error

		Stop func()

		ForceStop  func(i int)
		ForceIters int

		CancelContext bool
	}

	option func(o *settings)
)

var ErrKilled = errors.New("killed")

func Shutdown(ctx context.Context, run func(ctx context.Context) error, ops ...option) error {
	errc := make(chan error, 2)
	sigc := make(chan os.Signal, 1)

	s := settings{
		Signals:    []os.Signal{os.Interrupt},
		KillErr:    ErrKilled,
		ForceIters: 3,
	}

	for _, o := range ops {
		o(&s)
	}

	if !s.CancelContext && s.Stop == nil {
		return errors.New("no stop and no context cancel. how to stop it?")
	}

	signal.Notify(sigc, s.Signals...)

	var cancel func()

	if s.CancelContext {
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	go func() {
		errc <- run(ctx)
	}()

	go func() {
		<-sigc

		if cancel != nil {
			cancel()
		}

		if s.Stop != nil {
			s.Stop()
		}

		for i := 0; i < s.ForceIters; i++ {
			<-sigc

			if s.ForceStop != nil {
				s.ForceStop(i)
			}
		}

		errc <- s.KillErr
	}()

	return <-errc
}

func WithSignals(sig ...os.Signal) option {
	return func(o *settings) {
		o.Signals = sig
	}
}

func WithCancelContext() option {
	return func(o *settings) {
		o.CancelContext = true
	}
}

func WithStop(stop func()) option {
	return func(o *settings) {
		o.Stop = stop
	}
}

func WithForceStop(stop func(int)) option {
	return func(o *settings) {
		o.ForceStop = stop
	}
}

func WithForceIterations(n int) option {
	return func(o *settings) {
		o.ForceIters = n
	}
}
