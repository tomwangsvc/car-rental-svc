package os

//revive:disable:deep-exit

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

type Closer interface {
	Close()
}

type CloserWithError interface {
	Close() error
}

type Flusher interface {
	Flush()
}

func CleanUpAndExitOnInterrupt(ctx context.Context, closers []Closer, closerWithErrors []CloserWithError, flushers []Flusher) {
	lib_log.Info(ctx, "Initializing")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		signalReceived := <-c
		ctx := lib_context.NewCleanUpContext()

		lib_log.Notice(ctx, "Running clean up", lib_log.FmtString("signalReceived.String()", signalReceived.String()))
		for _, v := range closers {
			v.Close()
		}
		for i, v := range closerWithErrors {
			if err := v.Close(); err != nil {
				lib_log.Error(ctx, "Failed closing", lib_log.FmtInt("index", i), lib_log.FmtError(err))
			}
		}
		for _, v := range flushers {
			v.Flush()
		}
		lib_log.Notice(ctx, "Cleaned")

		os.Exit(0)
	}()

	lib_log.Info(ctx, "Initialized")
}
