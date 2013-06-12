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

func NewTypeServer(addr string) chan<- interface{} {
	typechan := make(chan interface{})
	go func(addr string, inchan <-chan interface{}) {
		server := typeServer{feed: inchan, write: handlers["html"]}
		err := http.ListenAndServe(addr, server)
		if err != nil {
			panic(err)
		}
	}(addr, typechan)
	return typechan
}

type typeWriter func(t *reflect.Type) (string, error)

var handlers = make(map[string]typeWriter)

type typeServer struct {
	feed  <-chan interface{}
	write typeWriter
}

func (server typeServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var t *reflect.Type
	if req.Method == "POST" {
		readType := reflect.TypeOf(<-server.feed)
		t = &readType
	}
	result, err := server.write(t)
	if err != nil {
		panic(err)
	}
	_, err = resp.Write([]byte(result))
	if err != nil {
		panic(err)
	}
}
