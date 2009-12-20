package xmlrpc

import("http"; "bytes"; "os")

// Calls a remote xmlrpc method. Response is the valid response form the server.
// Error may either be a local error or a remote one. If it is remote it will be of type
// xmlrpc.Fault.
func Call(endpoint string, req Request) (Response, os.Error) {
    body := new(bytes.Buffer)
    req.SendXML(body)
    resp, err := http.Post(endpoint, "text/xml", body)
    if err != nil {
        return Response{}, err
    }
    return ReadResponse(resp.Body)
}
