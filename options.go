package graceful

import (
	"context"
	"fmt"
	"path"
	"runtime"
)

type (
	Option interface {
		fmt.Stringer

		opt()
	}

	taskOpter interface {
		taskOpt(*task)
	}

	taskOpt struct {
		baseOpt

		f func(*task)
	}

	baseOpt string
)

func optFunc(d int) baseOpt {
	pc, _, _, ok := runtime.Caller(1 + d)
	if !ok {
		return baseOpt("<unknown>")
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return baseOpt("<unknown>")
	}

	name := f.Name()
	name = path.Ext(name)[1:]

	return baseOpt(name)
}

func (o baseOpt) String() string { return string(o) }

func WithStop(f func(context.Context) error) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			t.stop = f
		},
	}
}

func WithForceStop(f func(context.Context, int)) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			t.forceStop = f
		},
	}
}

func AllowStop(evenError bool) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			if evenError {
				t.allowStop = 2
			} else {
				t.allowStop = 1
			}
		},
	}
}

func WrapError(format string, args ...interface{}) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			t.wrapError = fmt.Sprintf(format, args...)
		},
	}
}

func IgnoreErrors(errs ...error) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			t.ignoreErrors = append(t.ignoreErrors, errs...)
		},
	}
}

func ErrorProcessor(f func(ctx context.Context, err error) error) Option {
	return taskOpt{
		baseOpt: optFunc(0),
		f: func(t *task) {
			t.processor = f
		},
	}
}

func (o taskOpt) taskOpt(t *task) {
	if o.f == nil {
		panic("not a task option")
	}

	o.f(t)
}

func (baseOpt) opt() {}
