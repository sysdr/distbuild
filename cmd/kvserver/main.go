// kvserver runs the distbuild key-value server.
//
//	go run ./cmd/kvserver -addr :7070
//	nc localhost 7070
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/systemdr/distbuild/internal/server"
	"github.com/systemdr/distbuild/internal/store"
)

func main() {
	addr := flag.String("addr", ":7070", "listen address")
	flag.Parse()

	srv := server.New(store.New())
	if err := srv.Listen(*addr); err != nil {
		log.Fatal(err)
	}
	log.Printf("kvserver listening on %s", srv.Addr())

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Print("shutting down")
		srv.Close()
	}()

	if err := srv.Serve(); err != nil {
		log.Fatal(err)
	}
}
