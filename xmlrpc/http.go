package xmlrpc

import("http"; "bytes"; "os")

func Call(endpoint string, req Request) (Response, os.Error) {
    body := new(bytes.Buffer)
    req.SendXML(body)
    resp, err := http.Post(endpoint, "text/xml", body)
    if err != nil {
        return Response{}, err
    }
    return ReadResponse(resp.Body)
}
