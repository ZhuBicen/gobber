// This file is part of Gobber, an open-source XMPP server written in Go
//
// Copyright © 2012 Josh Holland <jrh@joshh.co.uk>

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

// StreamError is a generic error related to a stream.
type StreamError struct {
	Condition string // one of the conditions from §4.9.3 of RFC 6120
	Message   string
}

func (s *StreamError) Error() string {
	return fmt.Sprintf("stream: %s: %s\n", s.Condition, s.Message)
}

// Stream represents an XMPP stream.
type Stream struct {
	io.ReadWriter

	To       string `xml:to,attr`
	From     string `xml:from,attr`
	Language string `xml:lang,attr`
	Id       string

	// Out is a channel that writes the given string response to the client
	Out chan string
}

// NewStream takes a ReadWriter and turns it into a stream
// if possible or returns an error otherwise.
func NewStream(buf io.ReadWriter) (s Stream, err error) {
	d := xml.NewDecoder(buf)
	s = Stream{ReadWriter: buf}

	s.Out = make(chan string)
	go func() {
		for val := range s.Out {
			io.WriteString(s, val)
		}
	}()

	for {
		t, err := d.Token()

		if err != nil {
			log.Print(err)
			break
		}

		switch el := t.(type) {
		case xml.StartElement:
			if el.Name.Local == "stream" {
				setupStream(&el, &s)
				return s, nil
			}
		}
	}

	return s, new(StreamError)
}

func setupStream(el *xml.StartElement, st *Stream) {
	for _, attr := range el.Attr {
		switch attr.Name.Local {
		case "to":
			st.To = attr.Value
		case "from":
			st.From = attr.Value
		case "lang":
			st.Language = attr.Value
		}
	}
}
