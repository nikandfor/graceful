package graceful

import (
	"context"
	"fmt"
	"path"

	"tlog.app/go/loc"
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
	pc := loc.Caller(1 + d)

	n, _, _ := pc.NameFileLine()

	n = path.Ext(n)[1:]

	return baseOpt(n)
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
