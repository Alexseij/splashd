//go:build linux

package unix

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var handleSignals = []os.Signal{
	unix.SIGTERM,
	unix.SIGINT,
}

func HandleSignals(
	log *logrus.Entry,
	doneCh chan struct{},
) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, handleSignals...)

	go func() {
		s := <-sigCh
		switch s {
		case unix.SIGTERM:
			log.Info("Catch SIGTERM")
		case unix.SIGINT:
			log.Info("Catch SIGINT")
		}
		doneCh <- struct{}{}
	}()
}
