package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/signal"
	"syscall"
)

func ImageBytesFmt(b []byte) string {
	_, imgType, err := image.Decode(bytes.NewBuffer(b))
	if err != nil {
		return "unknown"
	}

	return imgType
}

func WaitForSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	return quit
}
