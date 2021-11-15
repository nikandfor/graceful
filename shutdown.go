package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

type (
	opts struct {
		Signals []os.Signal

		KillErr error

		Stop func()

		ForceStop  func(i int)
		ForceIters int

		CancelContext bool
	}

	Option func(o *opts)
)

var ErrKilled = errors.New("killed")

func Shutdown(ctx context.Context, run func(ctx context.Context) error, ops ...Option) error {
	errc := make(chan error, 2)
	sigc := make(chan os.Signal, 1)

	s := opts{
		Signals:       []os.Signal{os.Interrupt},
		KillErr:       ErrKilled,
		ForceIters:    3,
		CancelContext: true,
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

		for i := s.ForceIters - 1; i >= 0; i-- {
			<-sigc

			if s.ForceStop != nil {
				s.ForceStop(i)
			}
		}

		<-sigc

		errc <- s.KillErr
	}()

	return <-errc
}

func WithSignals(sig ...os.Signal) Option {
	return func(o *opts) {
		o.Signals = sig
	}
}

func WithCancelContext(c bool) Option {
	return func(o *opts) {
		o.CancelContext = c
	}
}

func WithStop(stop func()) Option {
	return func(o *opts) {
		o.Stop = stop
	}
}

func WithForceStop(stop func(int)) Option {
	return func(o *opts) {
		o.ForceStop = stop
	}
}

func WithForceIterations(n int) Option {
	return func(o *opts) {
		o.ForceIters = n
	}
}
