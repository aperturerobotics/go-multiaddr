package main

import (
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

func main() {

	// construct from a string (err signals parse failure)
	m1, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	_ = err

	// construct from bytes (err signals parse failure)
	m2, err := ma.NewMultiaddrBytes(m1.Bytes())
	_ = err

	fmt.Printf("%v\n%v\n", m1, m2)

	listener, err := manet.Listen(m1)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("%v\n", listener)
}
