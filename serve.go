// typebrowser - View type information from your program in your browser!
//
// Copyright 2013 Arne Hormann and contributors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package typebrowser

import (
	"net/http"
	"reflect"
	"sync"
)

// Type is used to transport a message in addition to a type.
// Pass an instance to get a better overview which part of your code a type comes from.
type Type struct {
	Value   interface{}
	Message string
}

type convert func(message string, t *reflect.Type) (string, error)

type typeServer struct {
	inchan  <-chan interface{}
	mime    string
	convert convert
}

// make sure this is a http.Handler
var _ http.Handler = &typeServer{}

func (s *typeServer) NextString() string {
	data := <-s.inchan
	var t reflect.Type
	var m string
	if wrapped, ok := data.(Type); ok {
		t = reflect.TypeOf(wrapped.Value)
		m = wrapped.Message
	} else {
		t = reflect.TypeOf(data)
	}
	str, err := s.convert(m, &t)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *typeServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		// redirect to root for now
		resp.Header().Set("Location", "/")
		return
	}
	resp.Header().Set("Content-Type", s.mime)
	body := s.NextString()
	_, err := resp.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

type formServer struct{}

func (s formServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html")
	_, err := resp.Write([]byte(`<!DOCTYPE html><html><body>` + htmlForm + `</body></html>`))
	if err != nil {
		panic(err)
	}
}

var (
	typeConverters = make(map[string]typeConverter)
	setHtmlForm    sync.Once
	htmlForm       = ""
)

type typeConverter struct {
	mime    string
	convert convert
}

// NewTypeServer creates a new type server listening on the specified address.
// The address format matches that of http.ListenAndServe.
// You get a channel you can stuff things into.
// Point your browser to the address and fill the channel with values you want to inspect.
func NewTypeServer(addr string) chan<- interface{} {
	typechan := make(chan interface{})
	muxer := http.NewServeMux()
	form := ""
	for prefix, converter := range typeConverters {
		form += `<form method=post action="/` +
			prefix + `"><button type="submit">` +
			prefix + `</button></form>`
		muxer.Handle("/"+prefix, &typeServer{
			inchan:  typechan,
			mime:    converter.mime,
			convert: converter.convert,
		})
	}
	setHtmlForm.Do(func() {
		htmlForm = form
	})
	muxer.Handle("/", formServer{})
	go func(addr string, handler http.Handler) {
		if err := http.ListenAndServe(addr, handler); err != nil {
			panic(err)
		}
	}(addr, muxer)
	return typechan
}
