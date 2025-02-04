package graceful_test

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"nikand.dev/go/graceful"
)

func ExampleGroup() {
	httpAddr := flag.String("http", "", "http server address to listen to")

	flag.Parse()

	g := graceful.New()

	if *httpAddr != "" {
		l, err := net.Listen("tcp", *httpAddr)
		_ = err // if err != nil ...

		defer l.Close()

		g.Add(func(ctx context.Context) error {
			return http.Serve(l, nil)
		},
			graceful.WrapError("http server"),
			graceful.WithStop(func(ctx context.Context) error { return l.Close() }),
		)
	}

	g.Add(func(ctx context.Context) error {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()

		for {
			select {
			case <-t.C:
			case <-ctx.Done():
				return ctx.Err()
			}

			log.Printf("30 seconds has passed")
		}
	},
		graceful.WrapError("reporter"),
		graceful.IgnoreErrors(context.Canceled),
	)

	err := g.Run(context.Background())
	_ = err // process error
}
