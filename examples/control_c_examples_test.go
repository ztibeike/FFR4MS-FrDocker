package examples

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestControlC(t *testing.T) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)

}
