package xmlrpc

import (
	"http"
	"bytes"
	"os"
)

// Calls a remote xmlrpc method. Response is the valid response form the server.
// Error may either be a local error or a remote one. If it is remote it will be of type
// xmlrpc.Fault.
func (r RemoteMethod) Call(args ...) (Response, os.Error) {
	body := new(bytes.Buffer)
	r.SendXML(body, Params(args))
	resp, err := http.Post(r.Endpoint, "text/xml", body)
	if err != nil {
		return Response{}, err
	}
	return ReadResponse(resp.Body)
}
