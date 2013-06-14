typebrowser
===========

View type information from your program in your browser!

Typebrowser is a modern variant of println-debugging.
Instead of debug outputs on some lines scrolling faster than you can say ... anything,
it conviniently stops everything to let you see what's happening inside.

For now, this is restricted to type information.

How do I use this?
------------------
Call `NewTypeServer` and stuff things into the channel.
Open your browser and point it to the port you specify.
It's pretty simple.
```go
import (
	// import it
	"github.com/arnehormann/typebrowser"
)

// start it
var typechan = typebrowser.NewTypeServer(":8080")

...
	// use it - feed anything to it!
	typechan <- myvar
	// if you want to know where in your program this is called from:
	// declare a new type, assign your variable and pass state information as
	// an inspectable type tag
	typechan <- struct{value: interface{} `some file, some line, some state`}{myvar}
...
```
Typebrowser can be installed by `go get github.com/arnehormann/typebrowser`.

It depends on mirror, go-gettable with `go get github.com/arnehormann/mirror`.

When you use it, it shows a site with one button per export format.
The next type is fetched with an http-post, your program stops because it blocks on the channel read until the next type is requested.

For a working code example, see [example/example.go](example/example.go).

For a peek at what the html output looks like, see [this demo for *reflect.StructField](http://bl.ocks.org/arnehormann/raw/5780257/).

Types containing others are foldable on click.

JSON output is next on my agenda.

License: [MPL2](https://github.com/arnehormann/typebrowser/blob/master/LICENSE.md).
