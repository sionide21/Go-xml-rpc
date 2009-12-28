// Copyright 2009 The Ben Olive. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package xmlrpc

import (
	"reflect"
	"fmt"
)

// Given a list of args return an array of `ParamValue`s
func Params(params ...) []ParamValue {
	pStruct := reflect.NewValue(params).(*reflect.StructValue)
	par := make([]ParamValue, pStruct.NumField())

	for n := 0; n < len(par); n++ {
		par[n] = param(pStruct.Field(n))
	}
	return par
}

func argLength(args ...) int {
	pStruct := reflect.NewValue(args).(*reflect.StructValue)
	return pStruct.NumField()
}

func structParams(v *reflect.StructValue) StructValue {
	p := make(StructValue, v.NumField())
	for n := 0; n < v.NumField(); n++ {
		f := v.Type().(*reflect.StructType).Field(n)
		key := f.Name
		if f.Tag != "" {
			key = f.Tag
		}
		p[key] = param(v.Field(n))
	}
	return p
}

func mapParams(v *reflect.MapValue) StructValue {
	p := make(StructValue, v.Len())
	for _, k := range v.Keys() {
		key := k.(*reflect.StringValue).Get()
		p[key] = param(v.Elem(k))
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

func arrayParams(params *reflect.SliceValue) ArrayValue {
	par := make([]ParamValue, params.Len())

	for n := 0; n < len(par); n++ {
		par[n] = param(params.Elem(n))
	}
	return par
}

func param(param interface{}) ParamValue {
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
		return arrayParams(v)
	case *reflect.StructValue:
		return structParams(v)
    case *reflect.MapValue:
        // Might be a map[string]ParamValue which is treated as a struct
        ty, _ := v.Type().(*reflect.MapType)
        if _, ok := ty.Key().(*reflect.StringType); ok {
            return mapParams(v)
        }
	case *reflect.InterfaceValue:
		// If it is already a param value just return it
		if ret, ok := v.Interface().(ParamValue); ok {
			return ret
		}
	}
	// TODO How should this error be handled?
	return StringValue(fmt.Sprintf("Error: Unknown Param Type (%T)\n", param))
}
