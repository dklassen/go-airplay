package airplay

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/oleksandr/bonjour"
)

// Debug logger - main can reach in an enable this if it wants
var Debug = log.New(ioutil.Discard, "DEBUG ", log.LstdFlags)

// Default TXT Record, have not come up with a new one yet.
var txt map[string]string = map[string]string{
	"txtvers": "1",
	"pw":      "false",
	"tp":      "UDP",
	"sm":      "false",
	"ek":      "1",
	"cn":      "0,1",
	"ch":      "2",
	"ss":      "16",
	"sr":      "44100",
	"vn":      "3",
	"et":      "0,1",
}

func RegisterAirTunes(name, address string) (*bonjour.Server, error) {
	ifacename := "en0"

	_, portstr, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, err
	}
	iface, err := net.InterfaceByName(ifacename)
	if err != nil {
		return nil, err
	}

	raopName := hex.EncodeToString(iface.HardwareAddr) + "@" + name
	keys := make([]string, 0, len(txt))
	for key, value := range txt {
		keys = append(keys, fmt.Sprintf("%s=%s", key, value))
	}

	s, err := bonjour.Register(raopName, "_raop._tcp", "", port, keys, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return s, nil
}

func ServeAirTunes(name, address string, listener net.Listener, handler func(string, net.Conn)) error {
	log.Println("Listening for connections on address", address)

	for {
		id := strconv.Itoa(rand.Int())
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handler(id, conn)
	}
}
