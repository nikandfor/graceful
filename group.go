package graceful

import (
	"context"
	"os"
	"os/signal"

	"github.com/nikandfor/errors"
)

type (
	Group struct {
		Signals []os.Signal

		ForceIters int

		KillErr    error
		NoTasksErr error

		tasks []task
	}

	task struct {
		name string
		ctx  context.Context

		// set by user
		run       func(ctx context.Context) error
		stop      func(ctx context.Context) error
		forceStop func(ctx context.Context, i int)

		// context
		cancel func()

		allowStop bool

		done chan struct{}
	}
)

var (
	ErrKilled  = errors.New("killed")
	ErrNoTasks = errors.New("no tasks")
)

func New() *Group {
	return &Group{
		Signals:    []os.Signal{os.Interrupt},
		KillErr:    ErrKilled,
		NoTasksErr: ErrNoTasks,
		ForceIters: 3,
	}
}

func (g *Group) Add(ctx context.Context, name string, run func(context.Context) error, opts ...Option) {
	t := task{
		name: name,
		ctx:  ctx,
		run:  run,
		done: make(chan struct{}),
	}

	for _, o := range opts {
		if o, ok := o.(taskOpter); ok {
			o.taskOpt(&t)
		}
	}

	g.tasks = append(g.tasks, t)
}

func (g *Group) Run(ctx context.Context, opts ...Option) (err error) {
	if len(g.tasks) == 0 {
		return g.NoTasksErr
	}

	errc := make(chan error, 1)

	for _, t := range g.tasks {
		t := t
		go func() {
			defer close(t.done)

			err := t.run(t.ctx)

			if t.name != "" {
				err = errors.Wrap(err, t.name)
			}

			select {
			case errc <- err:
			default: // only first error matters
			}
		}()
	}

	var killc chan struct{}
	if g.Signals != nil {
		fin := make(chan struct{})
		defer close(fin) // to allow goroutine to finish

		killc = g.sigTask(errc, fin)
	}

	for _, t := range g.tasks {
		select {
		case <-t.done:
		case <-killc:
		}
	}

	return <-errc
}

func (g *Group) sigTask(errc chan error, fin chan struct{}) (kill chan struct{}) {
	sigc := make(chan os.Signal, 1)
	kill = make(chan struct{})

	signal.Notify(sigc, g.Signals...)

	go func() {
		select {
		case <-sigc:
		case <-fin:
			return
		}

		err := g.stop()

		if err != nil {
			select {
			case errc <- err:
			default: // only first error matters
			}
		}

		for i := g.ForceIters - 1; i >= 0; i-- {
			select {
			case <-sigc:
			case <-fin:
				return
			}

			g.forceStop(i)
		}

		select {
		case <-sigc:
		case <-fin:
			return
		}

		close(kill)

		select {
		case errc <- g.KillErr:
		case <-fin:
		}
	}()

	return
}

func (g *Group) stop() error {
	for _, t := range g.tasks {
		if t.cancel != nil {
			t.cancel()
		}

		if t.stop == nil {
			continue
		}

		err := t.stop(t.ctx)
		if err != nil {
			if t.name != "" {
				err = errors.Wrap(err, t.name)
			}

			return err
		}
	}

	return nil
}

func (g *Group) forceStop(i int) {
	for _, t := range g.tasks {
		if t.forceStop != nil {
			t.forceStop(t.ctx, i)
		}
	}
}
