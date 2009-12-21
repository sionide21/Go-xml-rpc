// Copyright 2009 The Ben Olive. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"io"
	"xml"
	"strings"
	"os"
)

type Response struct {
	Value ParamValue
}

// This methods reads the xml representation of a response
// to a io.Writer. It is public mostly for debugging or alternate
// communication. To call a remote function you should use
// `xmlrpc.Call`.
func ReadResponse(in io.Reader) (Response, os.Error) {
	p := tokenStream{xml.NewParser(in)}
	p.next(false) // Discard <methodResponse>
	t, err := p.next(false)
	if err != nil {
		return Response{}, err
	}
	resp, ok := t.(xml.StartElement)
	if !ok {
		return Response{}, error("Invalid Response")
	}
	switch strings.ToLower(resp.Name.Local) {
	case "fault":
		f, err := Fault{}.LoadXML(p.Parser)
		if err == nil {
			err = f.(Fault)
		}
		return Response{}, err
	case "params":
		p.next(false) // Eat <param>
		p.next(false) // Eat <value>
		a, err := parseMessage(p.Parser)
		return Response{a}, err
	}
	return Response{}, error("Invalid Response")
}

func (r Response) SendXML(out io.Writer) {
	io.WriteString(out, "<?xml version=\"1.0\"?>\n")
	io.WriteString(out, "<methodResponse>")
	io.WriteString(out, "<params>")
	io.WriteString(out, "<param><value>")
	io.WriteString(out, r.Value.ToXML())
	io.WriteString(out, "</value></param>")
	io.WriteString(out, "</params>")
	io.WriteString(out, "</methodResponse>")
}

func (f Fault) SendXML(out io.Writer) {
	io.WriteString(out, "<?xml version=\"1.0\"?>\n")
	io.WriteString(out, "<methodResponse>")
	io.WriteString(out, "<fault>")
	io.WriteString(out, "<value>")
	io.WriteString(out, f.ToXML())
	io.WriteString(out, "</value>")
	io.WriteString(out, "</fault>")
	io.WriteString(out, "</methodResponse>")
}
