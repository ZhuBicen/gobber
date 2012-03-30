// This file is part of Gobber, an open-source XMPP server written in Go
//
// Copyright © 2012 Josh Holland <jrh@joshh.co.uk>

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
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

	To, From string
	Language string
	Id       string
}

// NewStream takes a ReadWriter and turns it into a stream
// if possible or returns an error otherwise.
// It verifies the xml header and initial stream element and
// sets up the stream object.
func NewStream(buf io.ReadWriter) (s Stream, err error) {
	s = Stream{ReadWriter: buf}

	d := xml.NewDecoder(buf)

	header, err := getProcInst(d)

	if err != nil {
		return s, err
	}

	if header.Target != "xml" {
		return s, &StreamError{"bad-format", "ProcInst not directed at xml"}
	}

	// consume whitespace
	data, err := getCharData(d)
	if err == nil && len(bytes.TrimSpace(data)) > 0 {
		return s, &StreamError{"bad-format", "Invalid characters"}
	}

	elem, err := getStartElement(d)

	if n := elem.Name; n.Local != "stream" || n.Space != "stream" {
		return s, &StreamError{"bad-format", "Start element should be stream:stream"}
	}

	for _, attr := range elem.Attr {
		switch attr.Name.Local {
		case "to":
			s.To = attr.Value
		case "from":
			s.From = attr.Value
		case "id":
			s.Id = attr.Value
		case "lang":
			s.Language = attr.Value
		case "version":
			if !CheckVersion(attr.Value) {
				return s, &StreamError{"unsupported-version", attr.Value}
			}
		case "stream":
			if space, val := attr.Name.Space, attr.Value; space != "xmlns" || val != "http://etherx.jabber.org/streams" {
				return s, &StreamError{"invalid-namespace", attr.Value}
			}
		}
	}

	return s, nil
}

// CheckVersion checks to see whether we can handle this version
// of XMPP, as per RFC 6120 §4.7.5. It currently checks only that the
// version is correctly formatted and that the major part is 1. It will
// return true if we can handle it and false if not.
func CheckVersion(version string) bool {
	fields := strings.Split(version, ".")
	if len(fields) != 2 {
		return false
	}

	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return false
	}

	_, err = strconv.Atoi(fields[1])
	if err != nil {
		return false
	}

	if major != 1 {
		return false
	}

	return true
}
