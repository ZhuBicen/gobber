// This file is part of Gobber, an open-source XMPP server written in Go
//
// Copyright © 2012 Josh Holland <jrh@joshh.co.uk>

package main

import (
	"encoding/xml"
	"fmt"
)

// WrongToken signifies that an unexpected token type was
// next in the stream.
type WrongToken struct {
	Wanted, Got string
}

func (s *WrongToken) Error() string {
	return fmt.Sprintf("WrongToken: expected %s and got %s", s.Wanted, s.Got)
}

func getStartElement(d *xml.Decoder) (elem xml.StartElement, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.StartElement)
	if !ok {
		return elem, &WrongToken{"xml.StartElement", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}

func getEndElement(d *xml.Decoder) (elem xml.EndElement, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.EndElement)
	if !ok {
		return elem, &WrongToken{"xml.EndElement", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}

func getCharData(d *xml.Decoder) (elem xml.CharData, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.CharData)
	if !ok {
		return elem, &WrongToken{"xml.CharData", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}

func getComment(d *xml.Decoder) (elem xml.Comment, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.Comment)
	if !ok {
		return elem, &WrongToken{"xml.Comment", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}

func getProcInst(d *xml.Decoder) (elem xml.ProcInst, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.ProcInst)
	if !ok {
		return elem, &WrongToken{"xml.ProcInst", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}

func getDirective(d *xml.Decoder) (elem xml.Directive, err error) {
	token, err := d.RawToken()
	if err != nil {
		return
	}

	elem, ok := token.(xml.Directive)
	if !ok {
		return elem, &WrongToken{"xml.Directive", fmt.Sprintf("%T", token)}
	}

	return elem, nil
}
