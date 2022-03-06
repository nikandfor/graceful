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
		mctx *multicontext

		// set by user
		run       func(ctx context.Context) error
		stop      func(ctx context.Context) error
		forceStop func(ctx context.Context, i int)

		// context
		cancel func()

		allowStop int // 0 - don't, 1 - allow with nil error, 2 - allow with error
		//	stopCb    func(ctx context.Context, err error)

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

	if t.cancel == nil {
		t.ctx, t.cancel = context.WithCancel(t.ctx)
	}

	g.tasks = append(g.tasks, t)
}

/*
Plan:
    * Run all tasts concurrently
    * Wait for the first to finish
    * Stop all other tasts (cancel context)
    * Wait for all tasts to finish
    * Kill if not finished
	* Return the first non-nil error (or nil)

If one of Group.Signals is received all tasks are stopped.
If Group.ForceIters more signals received Group.Run returns immediately.
*/
func (g *Group) Run(ctx context.Context, opts ...Option) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	errc := make(chan error, len(g.tasks))

	for _, t := range g.tasks {
		t := t

		go func() {
			defer close(t.done)

			ctx := multi(t.ctx, ctx)

			err := t.run(ctx)

			//	if t.stopCb != nil {
			//		t.stopCb(ctx, err)
			//	}

			if t.allowStop > 1 || t.allowStop > 0 && err == nil {
				return
			}

			err = errors.Wrap(err, t.name)

			errc <- err
		}()
	}

	var sigc chan os.Signal
	if len(g.Signals) != 0 {
		sigc = make(chan os.Signal, 1)

		signal.Notify(sigc, g.Signals...)
		defer signal.Stop(sigc)
	}

	select {
	case err = <-errc:
	case <-ctx.Done():
		err = ctx.Err()
	case <-sigc:
	}

	e := g.stop()
	if err == nil {
		err = e
	}

	toKill := g.ForceIters

next:
	for _, t := range g.tasks {
		for {
			select {
			case <-t.done:
				continue next
			case <-sigc:
			}

			if toKill == 0 {
				continue next
			}

			toKill--

			g.forceStop(toKill)
		}
	}

	for err == nil && len(errc) != 0 {
		err = <-errc
	}

	return err
}

func (g *Group) stop() (err error) {
	for _, t := range g.tasks {
		t.cancel()

		if t.stop == nil {
			continue
		}

		e := t.stop(t.mctx)
		if err == nil {
			err = errors.Wrap(e, "stop: %v", t.name)
		}
	}

	return err
}

func (g *Group) forceStop(i int) {
	for _, t := range g.tasks {
		select {
		case <-t.done:
			continue
		default:
		}

		if t.forceStop != nil {
			t.forceStop(t.mctx, i)
		}
	}
}
