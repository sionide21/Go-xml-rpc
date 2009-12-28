// Copyright 2009 The Ben Olive. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"http"
	"bytes"
	"os"
	"time"
	"fmt"
)

// Calls a remote xmlrpc method. Response is the valid response form the server.
// Error may either be a local error or a remote one. If it is remote it will be of type
// xmlrpc.Fault.
func (r RemoteMethod) Call(args ...) (ParamValue, os.Error) {
	body := new(bytes.Buffer)
	if r.BaseParams != nil {
		p := make([]ParamValue, len(r.BaseParams)+argLength(args))
		i := 0
		for _, v := range r.BaseParams {
			p[i] = v
			i++
		}
		for _, v := range Params(args) {
			p[i] = v
			i++
		}
		r.SendXML(body, p)
	} else {
		r.SendXML(body, Params(args))
	}
	resp, err := http.Post(r.Endpoint, "text/xml", body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, error(fmt.Sprintf("%s (%d)", resp.Status, resp.StatusCode))
	}
	ret, err := ReadResponse(resp.Body)
	return ret.Value, err
}

func (r RemoteMethod) CallInt(args ...) (int, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return 0, err
	}
	return int(res.(IntValue)), nil
}

func (r RemoteMethod) CallBoolean(args ...) (bool, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return false, err
	}
	return bool(res.(BooleanValue)), nil
}

func (r RemoteMethod) CallDouble(args ...) (float, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return 0, err
	}
	return float(res.(DoubleValue)), nil
}

func (r RemoteMethod) CallDate(args ...) (time.Time, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return time.Time{}, err
	}
	return time.Time(res.(DateTimeValue)), nil
}

func (r RemoteMethod) CallBytes(args ...) ([]byte, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return nil, err
	}
	return []byte(res.(Base64Value)), nil
}

func (r RemoteMethod) CallStruct(args ...) (StructValue, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return nil, err
	}
	return res.(StructValue), nil
}

func (r RemoteMethod) CallArray(args ...) (ArrayValue, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return nil, err
	}
	return res.(ArrayValue), nil
}

func (r RemoteMethod) CallString(args ...) (string, os.Error) {
	res, err := r.Call(args)
	if err != nil {
		return "", err
	}
	return string(res.(StringValue)), nil
}
