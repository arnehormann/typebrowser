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
)

type typeWriter func(t *reflect.Type) (string, error)

type chanSourcer struct {
	inchan <-chan interface{}
	write  func(t *reflect.Type) (string, error)
}

// make sure this is a http.Handler
var _ http.Handler = &chanSourcer{}

func (s *chanSourcer) NextString() string {
	readType := reflect.TypeOf(<-s.inchan)
	str, err := s.write(&readType)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *chanSourcer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var body string
	var err error
	if req.Method == "POST" {
		body = s.NextString()
	} else {
		body, err = s.write(nil)
		if err != nil {
			panic(err)
		}
	}
	_, err = resp.Write([]byte(body))
	if err != nil {
		panic(err)
	}
}

var typeWriters = make(map[string]typeWriter)

func NewTypeServer(addr string) chan<- interface{} {
	typechan := make(chan interface{})
	muxer := http.NewServeMux()
	for prefix, write := range typeWriters {
		muxer.Handle("/"+prefix, &chanSourcer{inchan: typechan, write: write})
	}
	go func(addr string, handler http.Handler) {
		if err := http.ListenAndServe(addr, handler); err != nil {
			panic(err)
		}
	}(addr, muxer)
	return typechan
}
