package graceful

import (
	"context"
	"sync"
	"time"
)

type (
	multicontext struct {
		a, b context.Context

		doneOnce sync.Once
		done     <-chan struct{}
		errf     func() error

		multivalue bool
	}
)

var _ context.Context = &multicontext{}

func Multicontext(a, b context.Context) context.Context {
	c := multi(a, b)
	c.multivalue = true

	return c
}

func Multicancel(a, b context.Context) context.Context {
	return multi(a, b)
}

func multi(a, b context.Context) *multicontext {
	if a == nil {
		panic("nil a parent")
	}
	if b == nil {
		panic("nil b parent")
	}

	c := &multicontext{
		a: a,
		b: b,
	}

	return c
}

func (c *multicontext) Deadline() (d time.Time, ok bool) {
	d, ok = c.a.Deadline()

	d2, ok2 := c.b.Deadline()

	if !ok || ok2 && d2.Before(d) {
		return d2, ok2
	}

	return d, ok
}

func (c *multicontext) Done() <-chan struct{} {
	c.doneOnce.Do(func() {
		ad := c.a.Done()
		bd := c.b.Done()

		if ad == nil && bd == nil {
			return
		}

		if ad == nil {
			c.done = bd
			c.errf = c.b.Err

			return
		}

		if bd == nil {
			c.done = ad
			c.errf = c.a.Err

			return
		}

		done := make(chan struct{})
		c.done = done

		go func() {
			defer close(done)

			select {
			case <-ad:
				c.errf = c.a.Err
			case <-bd:
				c.errf = c.b.Err
			}
		}()
	})

	return c.done
}

func (c *multicontext) Err() error {
	select {
	case <-c.Done():
		return c.errf()
	default:
		return nil
	}
}

func (c *multicontext) Value(k interface{}) (v interface{}) {
	v = c.a.Value(k)

	if c.multivalue && v == nil {
		v = c.b.Value(k)
	}

	return
}
