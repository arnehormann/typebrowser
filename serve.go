// typebrowser - View type information from your program in your browser!
//
// Copyright 2013 Arne Hormann and contributors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package typebrowser

import (
	"bufio"
	"fmt"
	"github.com/arnehormann/mirror"
	"net/http"
	"reflect"
	"strings"
)

func NewTypeServer(addr string) chan<- interface{} {
	typechan := make(chan interface{})
	go func(addr string, inchan <-chan interface{}) {
		server := typeServer{feed: inchan, write: htmlTypeWriter}
		err := http.ListenAndServe(addr, server)
		if err != nil {
			panic(err)
		}
	}(addr, typechan)
	return typechan
}

type typeWriter func(s *typeSession, t *reflect.Type) error

type typeServer struct {
	feed  <-chan interface{}
	write typeWriter
}

type typeSession struct {
	depth int
	buf   *bufio.Writer
	err   error
}

func (server typeServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session := &typeSession{
		depth: 0,
		buf:   bufio.NewWriter(resp),
	}
	var t *reflect.Type
	if req.Method == "POST" {
		readType := reflect.TypeOf(<-server.feed)
		t = &readType
	}
	err := server.write(session, t)
	if err != nil {
		panic(err)
	}
	session.buf.Flush()
}

func (session *typeSession) Concat(text string) {
	if session.err != nil {
		return
	}
	_, err := session.buf.WriteString(text)
	session.err = err
}

func (session *typeSession) Concatf(format string, args ...interface{}) {
	if session.err != nil {
		return
	}
	_, err := session.buf.WriteString(fmt.Sprintf(format, args...))
	session.err = err
}

// code for html type export

func htmlTypeWriter(session *typeSession, t *reflect.Type) error {
	const submit = `<form method="post"><button type="submit">Next</button></form>`
	if t == nil {
		// serve form on GET requests so favicon.ico and co don't skip object under inspection
		session.Concat(`<!DOCTYPE html><html><body>` + submit + `</body></html>`)
		return session.err
	}
	// write leading...
	session.Concatf(`
<!DOCTYPE html>
<html><head><title>Go: '%s'</title><style>
html { background-color: #fafafa; }
div[data-kind] {
	box-sizing: border-box;
	position: relative;
	/* font */
	font-family: "HelveticaNeue-Light", "Helvetica Neue Light", "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
	font-weight: 300;
	font-size: 16px;
	line-height: 1.5em;
	color: #444444;
	/* defaults */
	border: none;
	border-color: #eeeeee;
	border-left: 1.5em solid;
	border-top: 4px solid;
	padding: 0.5em 0 0 0.5em;
}
div[data-kind]::before {
	content: attr(data-kind) ': ' attr(data-field) ' ' attr(data-type);
	position: relative;
	margin-left: 1em;
}
div[data-kind=int8],
div[data-kind=int16],
div[data-kind=int32],
div[data-kind=int64],
div[data-kind=int]				{ border-color: #0f808c; }

div[data-kind=uint8],
div[data-kind=uint16],
div[data-kind=uint32],
div[data-kind=uint64],
div[data-kind=uint]				{ border-color: #198c6f; }

div[data-kind=float32],
div[data-kind=float64]			{ border-color: #5b8c39; }

div[data-kind=complex64],
div[data-kind=complex128]		{ border-color: #778c1b; }

div[data-kind=bool]				{ border-color: #19758c; }
div[data-kind=rune]				{ border-color: #4e398c; }
div[data-kind=ptr]				{ border-color: #d96485; }

div[data-kind=uintptr],
div[data-kind="unsafe.Pointer"]	{ border-color: #d91d29; }

div[data-kind=array],
div[data-kind=slice]			{ border-color: #f29a19; }

div[data-kind=string]			{ border-color: #40478c; }
div[data-kind=map]				{ border-color: #f2C91f; }
div[data-kind=struct]			{ border-color: #8Ab048; }
div[data-kind=chan]				{ border-color: #9c0c40; }
div[data-kind=interface]		{ border-color: #5d277d; }
div[data-kind=func]				{ border-color: #7d0a72; }

.parent { color: red;  cursor: pointer; }
.hide * { display: none; }
.parent.hide::after {
	color: blue;
	content: ' [+]';
}
</style>
</head><body>%s`, *t, submit)
	typeToHtml := func(t *reflect.StructField, typeIndex, depth int) error {
		// close open tags
		if session.depth > depth {
			session.Concat(strings.Repeat("</div>", session.depth-depth))
		}
		// close this tag later
		session.depth = depth + 1
		// if no type is given, return
		if t == nil {
			return nil
		}

		classes := ""

		tt := t.Type
		session.Concatf(
			`<div data-kind="%s" data-type="%s" data-size="%d" data-typeid="%d"`,
			tt.Kind(), tt, tt.Size(), typeIndex)
		if len(t.Index) > 0 {
			session.Concatf(
				` data-field="%s" data-index="%v" data-offset="%d" data-tag="%s"`,
				t.Name, t.Index, t.Offset, t.Tag)
		}
		switch tt.Kind() {
		case reflect.Array:
			session.Concatf(` data-length="%d"`, tt.Len())
			classes += "parent "
		case reflect.Chan:
			var direction string
			switch tt.ChanDir() {
			case reflect.RecvDir:
				direction = "receive"
			case reflect.SendDir:
				direction = "send"
			case reflect.BothDir:
				direction = "both"
			}
			session.Concat(` data-direction="` + direction + `"`)
			classes += "parent "
		case reflect.Map:
			session.Concatf(` data-keytype="%s"`, tt.Key())
			classes += "parent "

		case reflect.Func:
			session.Concatf(` data-args-in="%d" data-args-out="%d"`, tt.NumIn(), tt.NumOut())

		case reflect.Ptr, reflect.Slice, reflect.Struct:
			classes += "parent "
		}
		session.Concat(` class="` + classes + `">`)
		return session.err
	}
	// walk the type
	session.err = mirror.Walk(*t, typeToHtml)
	if session.err != nil {
		return session.err
	}
	// close all tags
	session.err = typeToHtml(nil, 0, 0)
	// write closing code...
	session.Concat(`
<script>
var parents = document.getElementsByClassName('parent');
for(var i = 0; i < parents.length; i++) {
    parents[i].onclick = function(e) {
    	e.stopPropagation();

    	// toggle class 'hide'
    	this.className = (this.classList.contains('hide')) ?
    		this.className.replace(/hide/,'') :
    		this.className += ' hide';
    }
}
</script></body></html>`)
	return session.err
}
