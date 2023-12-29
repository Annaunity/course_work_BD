package shutdown

import (
	"io"
	"os"
	"os/signal"
	"coursework/pkg/log"
)

// Hook needed to intercept the server shutdown signal
func Hook(signals []os.Signal, closeItems ...io.Closer) {
	logger := log.GetLogger()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc
	logger.Printf("caught signal %s. shutting down", sig)

	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Printf("failed to close %v: %v", closer, err)
		}
	}

	logger.Println("---server shutdown---")
}
