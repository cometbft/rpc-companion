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
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for sig := range c {
			logger.Info(fmt.Sprintf("Signal %v captured, exiting...", sig))
			if cb != nil {
				cb()
			}
			os.Exit(0)
		}
	}()
}
