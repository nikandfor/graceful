package graceful

import (
	"context"
	"os"
	"os/signal"

	"tlog.app/go/errors"
)

type (
	Group struct {
		Signals []os.Signal

		ForceIters int

		KillErr error
		//	NoTasksErr error

		tasks []task
	}

	task struct {
		//	name string

		// set by user
		run       func(ctx context.Context) error
		stop      func(ctx context.Context) error
		forceStop func(ctx context.Context, i int)
		allowStop int // 0 - don't, 1 - allow with nil error, 2 - allow with error

		wrapError    string
		ignoreErrors []error

		processor func(ctx context.Context, err error) error

		done chan struct{}
	}
)

var (
	ErrKilled  = errors.New("killed")
	ErrNoTasks = errors.New("no tasks")

	Restart = errors.New("restart")
)

func New() *Group {
	return &Group{
		Signals: []os.Signal{os.Interrupt},
		KillErr: ErrKilled,
		//	NoTasksErr: ErrNoTasks,
		ForceIters: 3,
	}
}

func Sub() *Group { return &Group{} }

func (g *Group) Add(run func(context.Context) error, opts ...Option) {
	t := task{
		run:  run,
		done: make(chan struct{}),
	}

	g.applyOpts(&t, opts)

	g.tasks = append(g.tasks, t)
}

/*
Plan:
  - Run all tasks concurrently
  - Wait for the first to finish
  - Stop all other tasks (cancel context)
  - Wait for all tasks to finish
  - Kill if not finished
  - Return the first non-nil error or nil

If one of Group.Signals is received all tasks are stopped.
If Group.ForceIters more signals received Group.Run returns immediately.
*/
func (g *Group) Run(ctx context.Context, opts ...Option) (err error) {
	select {
	case <-ctx.Done():
		return g.ctxErr(ctx, opts)
	default:
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errc := make(chan error, len(g.tasks))

	for i := range g.tasks {
		t := &g.tasks[i]

		g.applyOpts(t, opts)

		go g.taskRun(ctx, t, errc)
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
	case <-sigc:
	}

	cancel()

	e := g.stop(ctx)
	if err == nil {
		err = e
	}

	toKill := g.ForceIters

allTasks:
	for _, t := range g.tasks {
		for {
			select {
			case <-t.done:
				continue allTasks
			case <-sigc:
				toKill--
			}

			if toKill <= 0 {
				break allTasks
			}

			g.forceStop(ctx, toKill)
		}
	}

	for err == nil && len(errc) != 0 {
		err = <-errc
	}

	if err == nil && toKill <= 0 {
		err = g.KillErr
	}

	if err == nil {
		err = g.ctxErr(ctx, opts)
	}

	return err
}

func (g *Group) taskRun(ctx context.Context, t *task, errc chan error) {
	defer close(t.done)

restart:
	err := t.run(ctx)

	if t.processor != nil {
		err = t.processor(ctx, err)
		if errors.Is(err, Restart) {
			goto restart
		}
	}

	if t.allowStop > 1 || t.allowStop > 0 && err == nil {
		return
	}

	for _, ie := range t.ignoreErrors {
		if errors.Is(err, ie) {
			err = nil
			break
		}
	}

	if err != nil && t.wrapError != "" {
		err = errors.Wrap(err, t.wrapError)
	}

	errc <- err
}

func (g *Group) ctxErr(ctx context.Context, opts []Option) error {
	var t task
	g.applyOpts(&t, opts)

	err := ctx.Err()

	for _, ie := range t.ignoreErrors {
		if errors.Is(err, ie) {
			err = nil
			break
		}
	}

	if t.wrapError != "" {
		err = errors.Wrap(err, t.wrapError)
	}

	return err
}

func (g *Group) applyOpts(t *task, opts []Option) {
	for _, o := range opts {
		if o, ok := o.(taskOpter); ok {
			o.taskOpt(t)
		}
	}
}

func (g *Group) stop(ctx context.Context) (err error) {
	for _, t := range g.tasks {
		select {
		case <-t.done:
			continue
		default:
		}

		if t.stop == nil {
			continue
		}

		e := t.stop(ctx)
		if err == nil {
			err = e
			//	err = errors.Wrap(e, "stop: %v", t.name)
		}
	}

	return err
}

func (g *Group) forceStop(ctx context.Context, i int) {
	for _, t := range g.tasks {
		select {
		case <-t.done:
			continue
		default:
		}

		if t.forceStop != nil {
			t.forceStop(ctx, i)
		}
	}
}
