typebrowser
===========

View type information from your program in your browser!

How do I use this?
------------------
Call `NewTypeServer` and stuff things into the channel.
Open your browser and point it to the port you specify.
It's pretty simple.
```go
import (
	"github.com/arnehormann/typebrowser"
)

...
	typechan := typebrowser.NewTypeServer(":8080")
	typechan <- myvar
...
```
Typebrowser can be installed by `go get github.com/arnehormann/typebrowser`.
It depends on mirror, go-gettable with `go get github.com/arnehormann/mirror`.

When you use it, it shows a site with a button per export format. the next type is fetched with an http-post, your program stops because it blocks on the channel until the next type is requested.

For a working code example, see [example/example.go](example/example.go).
For a peek at what the html output looks like, see [this demo for *reflect.StructField](http://bl.ocks.org/arnehormann/raw/5775864/).
Types containing others are foldable on click.

For now, the JSON output is invalid. It's the next issue on my agenda.

License: [MPL2](https://github.com/arnehormann/typebrowser/blob/master/LICENSE.md).
