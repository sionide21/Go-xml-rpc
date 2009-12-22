// Copyright 2009 The Ben Olive. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"io"
	"fmt"
)

// The remote method to call. Endpoint is the URL of the endpoint
// Method is the name of the method to call.
type RemoteMethod struct {
	Endpoint       string
	Method         string
	RemoteFieldMap map[string](map[string]string)
}

func NewRemoteMethod(endpoint, method string) *RemoteMethod {
	return &RemoteMethod{endpoint, method, make(map[string](map[string]string))}
}

// This methods writes the xml representation of the request
// to a io.Writer. It is public mostly for debugging or alternate
// communication. To call a remote function you should use
// `xmlrpc.Call`.
func (r RemoteMethod) SendXML(out io.Writer, params []ParamValue) {
	io.WriteString(out, "<?xml version=\"1.0\"?>\n")
	io.WriteString(out, "<methodCall>")
	io.WriteString(out, fmt.Sprintf("<methodName>%s</methodName>", r.Method))
	io.WriteString(out, "<params>")
	for _, p := range params {
		io.WriteString(out, "<param><value>")
		io.WriteString(out, p.ToXML())
		io.WriteString(out, "</value></param>")
	}
	io.WriteString(out, "</params>")
	io.WriteString(out, "</methodCall>")
}
