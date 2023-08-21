package os

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// TrapSignal catches the SIGTERM/SIGINT and executes cb function. After that it exits
// with code 0.
func TrapSignal(logger slog.Logger, cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logger.Info("signal trapped", "msg", fmt.Sprintf("captured %v, exiting...", sig))
			if cb != nil {
				cb()
			}
			os.Exit(0)
		}
	}()
}
