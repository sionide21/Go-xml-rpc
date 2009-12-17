package xmlrpc

import ("io"; "fmt")

type Request struct {
    Method string
    Params []ParamValue
}

func (r Request) SendXML(out io.Writer) {
    io.WriteString(out, "<?xml version=\"1.0\"?>\n")
    io.WriteString(out, "<methodCall>")
    io.WriteString(out, fmt.Sprintf("<methodName>%s</methodName>", r.Method))
    io.WriteString(out, "<params>")
    for _, p := range r.Params {
        io.WriteString(out, "<param><value>")
        io.WriteString(out, p.ToXML());
        io.WriteString(out, "</value></param>")
    }
    io.WriteString(out, "</params>")
}
