// This file is part of Gobber, an open-source XMPP server written in Go
//
// Copyright © 2012 Josh Holland <jrh@joshh.co.uk>

package main

import (
	"log"
	"net"
)

func handle(conn net.Conn) {
	log.Printf("Handling connection from %s", conn.RemoteAddr())

	st, err := NewStream(conn)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got stream from %s to %s\n", st.From, st.To)
	log.Printf("Stream: %+v\n", st)
}
