// typebrowser - View type information from your program in your browser!
//
// Copyright 2013 Arne Hormann and contributors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package typebrowser

import (
	"fmt"
	"github.com/arnehormann/mirror"
	"reflect"
)

func init() {
	/* DISABLE for now
	typeConverters["json"] = typeConverter{
		mime:    `application/json`,
		convert: jsonConverter,
	}
	*/
}

func jsonConverter(message string, t *reflect.Type) (string, error) {
	// TODO: doesn't use message and format is broken
	if t == nil {
		return "{}", nil
	}
	lastDepth := 0
	opening := ""
	closing := ""
	var typeToJson func(t *reflect.StructField, typeIndex, depth int) error
	typeToJson = func(t *reflect.StructField, typeIndex, depth int) error {
		// close open tags
		if depthDelta := lastDepth - depth; depthDelta > 0 {
			opening += closing[:depthDelta]
			closing = closing[depthDelta:]
		}
		// close this tag later
		lastDepth = depth + 1
		// if no type is given, return
		if t == nil {
			return nil
		}
		tt := t.Type
		closing = "}" + closing
		opening += fmt.Sprintf(`{"kind":%q,"gotype":%q,"memsize":%d`, tt.Kind(), tt, tt.Size())
		if t.Name != "" {
			opening += fmt.Sprintf(`,"name":%q`, t.Name)
		}
		if len(t.Index) > 0 {
			if t.Tag != "" {
				opening += fmt.Sprintf(`,"tag":%q`, t.Tag)
			}
			opening += fmt.Sprintf(`,"structidx":%v,"memoffset":%d`, t.Index, t.Offset)
		}
		switch tt.Kind() {
		case reflect.Func:
			if tt.IsVariadic() {
				opening += `,"varargs":true`
			}
			opening += `,"arguments":[`
			for i, maxi := 0, tt.NumIn(); i < maxi; i++ {
				arg := tt.In(i)
				opening += fmt.Sprintf(`%q,`, arg.Name())
			}
			opening = opening[:len(opening)-1] + `],"returns":[`
			for i, maxi := 0, tt.NumOut(); i < maxi; i++ {
				ret := tt.Out(i)
				opening += fmt.Sprintf(`{"gotype":%q},`, ret.Name())
			}
			opening = opening[:len(opening)-1] + `]`
		case reflect.Interface:
			// iterate functions
			opening += fmt.Sprintf(`,"methods":[`)
			for i, maxi := 0, tt.NumMethod(); i < maxi; i++ {
				if i > 0 {
					opening += `,`
				}
				method := tt.Method(i)
				m := &reflect.StructField{Name: method.Name, Type: method.Type}
				err := typeToJson(m, -1, depth+1)
				if err != nil {
					return err
				}
			}
			opening += `]`
		case reflect.Struct:
			closing = "]" + closing
			opening += `,"fields":[`
		// ...
		case reflect.Chan:
			closing = "}" + closing
			switch tt.ChanDir() {
			case reflect.RecvDir:
				opening += `,"direction":"receive","elem"=`
			case reflect.SendDir:
				opening += `,"direction":"send","elem"=`
			case reflect.BothDir:
				opening += `,"direction":"both","elem"=`
			}
		case reflect.Map:
			closing = "}" + closing
			opening += fmt.Sprintf(`,"key":%q,"elem":`, tt.Key())
		case reflect.Invalid, reflect.Bool, reflect.Uintptr, reflect.UnsafePointer,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		}
		return nil
	}
	// walk the type
	err := mirror.Walk(*t, typeToJson)
	if err != nil {
		return "", err
	}
	return opening + closing, nil
}
