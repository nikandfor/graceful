# graceful

`graceful` is a library to start multiple routines and wait for them to finish gracefully. It handles errors and OS signals for you.

## Algorithm

* Add tasks to the group.
* Run all concurrently.
* If Group.Signals is not empty subscribe for signals.
* Wait for the first task to finish.
* Cancel context and call WithStop functions for each task if set.
* Wait for all tasks to finish.
* If more signals received call WithForceStop for each task if set.
* If Group.ForceIters signals received stop waiting.
* The first error occurred is returned.

## Uasage

```
g := graceful.New() // automatically handles interrupt signal while graceful.Sub doesn't.

g.Add(func(ctx context.Context) error {
    return s.Serve(ctx)
},
    graceful.WrapError("server"), // so we know where error come from if it fails
    graceful.WithStop(func(ctx context.Context) error {
        s.Shutdown()
        return nil
    }),
    graceful.WithForceStop(func(ctx context.Context, i int) error {
        return s.Kill()
    })
)

g.Add(func(ctx context.Context) error {
    return s.Worker(ctx)
},
    graceful.WrapError("worker"),
)

g.Add(func(ctx context.Context) error {
    return s.OneShotJob(ctx)
},
    graceful.WrapError("one shot job"),
    graceful.AllowStop(false), // if job exists, the rest of the group isn't stopped.
)

g.Add(func(ctx context.Context) error {
    return s.FailingWorker(ctx)
},
    graceful.WrapError("failing worker"),
    graceful.ErrorProcessor(func(ctx context.Context, err error) error {
        tlog.Printw("failing worker failed", "err", err)

        if unrecoverable(err) {
            return err
        }

        // backoff for some time

        return graceful.Restart
    })
)

return g.Run(ctx, gracefull.IgnoreErrors(context.Canceled))
```
