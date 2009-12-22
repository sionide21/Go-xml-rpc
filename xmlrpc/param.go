// Copyright 2009 The Ben Olive. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"reflect"
	"fmt"
)


// If you want the names of fileds in a struct to be different than on
// the wire than they are locally, specify the existing struct field name
// as a key and the desired name as the value.
// For Example, if I have the following struct:
//
// type A struct {
//   B string
// }
//
// But I wanted B to be called FileName for some remote method `remote`,
// I would call:
//
// r.RegisterFieldMap(A{}, map[staring]string{"B": "FileName"})
//
func (r RemoteMethod) RegisterFieldMap(s interface{}, mapping map[string]string) {
	typ := reflect.Typeof(s).String()
	if _, ok := r.RemoteFieldMap[typ]; !ok {
		r.RemoteFieldMap[typ] = make(map[string]string)
	}
	for k, v := range (mapping) {
		r.RemoteFieldMap[typ][k] = v
	}
}

// Given a list of args return an array of `ParamValue`s
func (r RemoteMethod) params(params ...) []ParamValue {
	pStruct := reflect.NewValue(params).(*reflect.StructValue)
	par := make([]ParamValue, pStruct.NumField())

	for n := 0; n < len(par); n++ {
		par[n] = param(r, pStruct.Field(n))
	}
	return par
}

func structParams(r RemoteMethod, v *reflect.StructValue) StructValue {
	p := make(StructValue, v.NumField())
	fieldMap, hasOverrides := r.RemoteFieldMap[v.Type().String()]
	for n := 0; n < v.NumField(); n++ {
		key := v.Type().(*reflect.StructType).Field(n).Name
		if hasOverrides {
			newKey, ok := fieldMap[key]
			if ok {
				key = newKey
			}
		}
		p[key] = param(r, v.Field(n))
	}
	return p
}

func byteParams(params *reflect.SliceValue) Base64Value {
	par := make([]byte, params.Len())

	for n := 0; n < len(par); n++ {
		par[n] = params.Elem(n).(*reflect.Uint8Value).Get()
	}
	return Base64Value(par)
}

func arrayParams(r RemoteMethod, params *reflect.SliceValue) ArrayValue {
	par := make([]ParamValue, params.Len())

	for n := 0; n < len(par); n++ {
		par[n] = param(r, params.Elem(n))
	}
	return par
}

func param(r RemoteMethod, param interface{}) ParamValue {
	switch v := param.(type) {
	case *reflect.IntValue:
		return IntValue(v.Get())
	case *reflect.BoolValue:
		return BooleanValue(v.Get())
	case *reflect.StringValue:
		return StringValue(v.Get())
	case *reflect.FloatValue:
		return DoubleValue(v.Get())
	case *reflect.SliceValue:
		if _, ok := v.Type().(*reflect.SliceType).Elem().(*reflect.Uint8Type); ok { // A []byte is really a Base64Type
			return byteParams(v)
		}
		return arrayParams(r, v)
	case *reflect.StructValue:
		return structParams(r, v)
	}
	// TODO How should this error be handled?
	return StringValue(fmt.Sprintf("Error: Unknown Param Type (%T)\n", param))
}
