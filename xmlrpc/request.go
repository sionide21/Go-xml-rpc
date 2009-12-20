package xmlrpc

import (
	"io"
	"fmt"
)

// The remote method to call. Method is the name of the method.
// Params is the list of arguments to pass to `Method`.
type Request struct {
	Method string
	Params []ParamValue
}

// This methods writes the xml representation of the request
// to a io.Writer. It is public mostly for debugging or alternate
// communication. To call a remote function you should use
// `xmlrpc.Call`.
func (r Request) SendXML(out io.Writer) {
	io.WriteString(out, "<?xml version=\"1.0\"?>\n")
	io.WriteString(out, "<methodCall>")
	io.WriteString(out, fmt.Sprintf("<methodName>%s</methodName>", r.Method))
	io.WriteString(out, "<params>")
	for _, p := range r.Params {
		io.WriteString(out, "<param><value>")
		io.WriteString(out, p.ToXML())
		io.WriteString(out, "</value></param>")
	}
	io.WriteString(out, "</params>")
	io.WriteString(out, "</methodCall>")
}
