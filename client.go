// This file is part of Gobber, an open-source XMPP server written in Go
//
// Copyright © 2012 Josh Holland <jrh@joshh.co.uk>

package main

import (
	"encoding/xml"
	"io"
	"log"
	"net"
)

// Client represents an XMPP client, encapsulating the two unidirectional
// streams to and from them.
type Client interface {
	// Close flushes the buffers and closes the stream to the client.
	Close() error

	// OutChannel and InChannel return channels for reading and
	// writing XMPP stanzas to the client respectively.
	InChannel() <-chan Stanza
	OutChannel() chan<- Stanza
}

func NewClient(conn net.Conn) (c Client, err error) {
	cl := client{Conn: conn}
	cl.Out = make(chan Stanza, 4)
	cl.In = make(chan Stanza, 4)

	// stuff

	go func() {
		for st := range cl.Out {
			io.Copy(cl.Conn, st)
		}
	}()

	go func() {
		d := xml.NewDecoder(conn)

		for {
			t, err := d.Token()

			if t == nil {
				log.Println(err)
				break
			}

			se, ok := t.(xml.StartElement)
			if !ok {
				log.Printf("Got %T, not xml.StartElement\n", t)
			}
			cl.In <- NewStanza(se)
		}
	}()
}

type client struct {
	net.Conn

	// Out is a channel to send stanzas to the client
	// As a convenience to avoid dealing with Write, send
	// stanzas to this channel to send them to the client
	Out chan Stanza

	// In is a channel to receive stanzas from the client
	// Like Out, this channel is for convenience
	In chan Stanza
}
