package graceful

import (
	"context"
)

func Shutdown(ctx context.Context, run func(ctx context.Context) error, opts ...Option) error {
	g := New()

	g.Add(ctx, "", run, opts...)

	return g.Run(ctx, opts...)
}
