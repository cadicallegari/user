package xsignal

import (
	"os"
	"os/signal"
	"syscall"
)

var StopSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func WaitSignal(sig ...os.Signal) os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	return <-c
}

func WaitStopSignal() os.Signal {
	return WaitSignal(StopSignals...)
}
