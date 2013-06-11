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

For a working example, see [example/example.go](example/example.go).

License: [MPL2](https://github.com/arnehormann/typebrowser/blob/master/LICENSE.md).
