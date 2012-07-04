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

	clientStream, err := NewStream(conn)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got stream from %s to %s\n", clientStream.From, clientStream.To)
	log.Printf("Stream: %+v\n", clientStream)

	log.Println("Replying to client")
	serverStream, err := NewStream(conn)
}
