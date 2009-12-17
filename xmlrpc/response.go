package xmlrpc

import ("io")

type Response struct {
    Value ParamValue
}

func (r Response) SendXML(out io.Writer) {
    io.WriteString(out, "<?xml version=\"1.0\"?>\n")
    io.WriteString(out, "<methodResponse>")
    io.WriteString(out, "<params>")
    io.WriteString(out, "<param><value>")
    io.WriteString(out, r.Value.ToXML());
    io.WriteString(out, "</value></param>")
    io.WriteString(out, "</params>")
    io.WriteString(out, "</methodResponse>")
}

func (f Fault) SendXML(out io.Writer) {
    io.WriteString(out, "<?xml version=\"1.0\"?>\n")
    io.WriteString(out, "<methodResponse>")
    io.WriteString(out, "<fault>")
    io.WriteString(out, "<value>")
    io.WriteString(out, f.ToXML());
    io.WriteString(out, "</value>")
    io.WriteString(out, "</fault>")
    io.WriteString(out, "</methodResponse>")
}

