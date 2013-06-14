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
	"html"
	"reflect"
	"strings"
)

func init() {
	typeConverters["html"] = typeConverter{
		mime:    `text/html`,
		convert: htmlConverter,
	}
}

func htmlConverter(message string, t *reflect.Type) (out string, err error) {
	if t == nil {
		return `<!DOCTYPE html><html></html>`, nil
	}
	lastDepth := 0
	Concat := func(text string) {
		out += text
	}
	Concatf := func(format string, args ...interface{}) {
		out += fmt.Sprintf(format, args...)
	}
	// write leading...
	Concatf(`<!DOCTYPE html>
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
	border-color: #eeeeee;
	padding: 0.5em 0 0 0.5em;
	/* enterprisify it a little */
	border: none;
	border-left: 1.5em solid;
	border-top: 4px solid;
	border-radius: 1em;
	border-top-right-radius: 0;
}
div[data-kind]::before {
	content: '[' attr(data-kind) ', ' attr(data-memsize) ' bytes]: ' attr(data-field) ' ' attr(data-type);
	position: relative;
	margin-left: 1em;
}
div[data-kind=int8]				{ border-color: hsl(180, 90%%, 50%%); }
div[data-kind=int16]			{ border-color: hsl(180, 90%%, 45%%); }
div[data-kind=int32]			{ border-color: hsl(180, 90%%, 40%%); }
div[data-kind=int64]			{ border-color: hsl(180, 90%%, 35%%); }
div[data-kind=int]				{ border-color: hsl(180, 75%%, 38%%); }
div[data-kind=uint8]			{ border-color: hsl(190, 90%%, 50%%); }
div[data-kind=uint16]			{ border-color: hsl(190, 90%%, 45%%); }
div[data-kind=uint32]			{ border-color: hsl(190, 90%%, 40%%); }
div[data-kind=uint64]			{ border-color: hsl(190, 90%%, 35%%); }
div[data-kind=uint]				{ border-color: hsl(190, 75%%, 38%%); }
div[data-kind=float32]			{ border-color: hsl(205, 70%%, 40%%); }
div[data-kind=float64]			{ border-color: hsl(205, 70%%, 35%%); }
div[data-kind=complex64]		{ border-color: hsl(215, 50%%, 35%%); }
div[data-kind=complex128]		{ border-color: hsl(215, 50%%, 30%%); }
div[data-kind=bool]				{ border-color: hsl(160, 70%%, 35%%); }
div[data-kind=ptr]				{ border-color: hsl(30, 50%%, 60%%); }
div[data-kind=uintptr]			{ border-color: hsl(20, 50%%, 50%%); }
div[data-kind="unsafe.Pointer"]	{ border-color: hsl(10, 90%%, 50%%); }
div[data-kind=array]			{ border-color: hsl(60, 90%%, 45%%); }
div[data-kind=slice]			{ border-color: hsl(60, 40%%, 60%%); }
div[data-kind=string]			{ border-color: hsl(120, 70%%, 30%%); }
div[data-kind=map]				{ border-color: hsl(75, 40%%, 40%%); }
div[data-kind=struct]			{ border-color: hsl(150, 10%%, 45%%); }
div[data-kind=interface]		{ border-color: hsl(240, 30%%, 60%%); }
div[data-kind=func]				{ border-color: hsl(270, 40%%, 60%%); }
div[data-kind=chan]				{ border-color: hsl(300, 40%%, 30%%); }

.fold * { display: none; }
.fold::after { content: ' [+]'; }
</style>
</head><body>`+htmlForm+`<hr>`, *t)
	if message != "" {
		Concatf("<h3>%s</h3><hr>\n", html.EscapeString(message))
	}
	expectInFunc := [][2]int{}
	typeToHtml := func(t *reflect.StructField, typeIndex, depth int) error {
		// close open tags
		if lastDepth > depth {
			Concat(strings.Repeat("</div>", lastDepth-depth))
		}
		// close this tag later
		lastDepth = depth + 1
		// if no type is given, return
		if t == nil {
			return nil
		}
		tt := t.Type
		Concatf(
			`<div data-kind="%s" data-type=%q data-memsize="%d" data-typeid="%d"`,
			tt.Kind(), html.EscapeString(tt.String()), tt.Size(), typeIndex)
		if len(expectInFunc) <= depth {
			expectInFunc = append(expectInFunc, [2]int{})
		} else {
			if expectInFunc[depth][0] > 0 {
				expectInFunc[depth][0]--
				Concat(` data-funcval="arg"`)
			} else {
				expectInFunc[depth][1]--
				Concat(` data-funcval="ret"`)
			}
		}
		if len(t.Index) > 0 {
			Concatf(
				` data-field="%s" data-index="%v" data-offset="%d" data-tag="%s"`,
				t.Name, t.Index, t.Offset, t.Tag)
		}
		switch tt.Kind() {
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
			Concat(` data-direction="` + direction + `"`)

		case reflect.Func:
			argcnt, retcnt := tt.NumIn(), tt.NumOut()
			if len(expectInFunc) <= depth+1 {
				expectInFunc = append(expectInFunc, [2]int{argcnt, retcnt})
			} else {
				expectInFunc[depth+1][0] = argcnt
				expectInFunc[depth+1][1] = retcnt
			}
			Concatf(` data-argcount="%d" data-retcount="%d"`, argcnt, retcnt)

		case reflect.Array:
			Concatf(` data-length="%d"`, tt.Len())

		case reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.Interface:
		}
		Concat(`>`)
		return nil
	}
	// walk the type
	err = mirror.Walk(*t, typeToHtml)
	if err != nil {
		return "", err
	}
	// close all tags
	_ = typeToHtml(nil, 0, 0)
	// write closing code...
	Concat(`
<script>
(function(tags, tag){
	function onChild(e) {
		e.stopPropagation()
	}
	function onParent(e) {
		e.stopPropagation()
		this.className = this.className == "fold" ? "" : "fold"
	}
	for (var i = 0; i < tags.length; i++) {
		tag = tags[i]
		tag.onclick = tag.children.length === 0 ? onChild : onParent
	}
})(document.getElementsByTagName('div'))
</script></body></html>`)
	return out, err
}
