// typebrowser - View type information from your program in your browser!
//
// Copyright 2013 Arne Hormann and contributors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package typebrowser

import (
	"reflect"
)

func init() {
	typeWriters[""] = formWriter
}

func formWriter(t *reflect.Type) (out string, err error) {
	return `<!DOCTYPE html><html><body>
Next as
<ul>
	<li><a href="/html">HTML</a>
	<li><a href="/json">JSON</a>
</ul>
</body></html>`, nil
}
