package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/joelgibson/go-airplay/airplay"
)

func main() {
	debug := flag.Bool("debug", false, "Show debug output.")
	name := flag.String("name", "AirSonos", "AirTunes name.")
	flag.Parse()

	if *debug {
		airplay.Debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
	}

	address := ":49152"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic("Listen failed: %s", err)
	}
	defer listener.Close()

	bonjourServer, err := airplay.RegisterAirTunes(*name, address)
	defer bonjourServer.Shutdown()

	err = airplay.ServeAirTunes(*name, address, listener, func(id string, conn net.Conn) {
		airplay.RtspSession(id, conn, func(x chan string) {})
	})

	if err != nil {
		log.Fatal(err)
	}
}
