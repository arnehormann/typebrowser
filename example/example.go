package main

import (
	"github.com/arnehormann/httpexport"
	"os/exec"
	"reflect"
	"runtime"
)

func main() {
	addr := ":8080"
	typechan := httpexport.NewTypeServer(addr)
	// open in browser
	cliOpener := "open"
	if runtime.GOOS == "windows" {
		cliOpener = "start"
	}
	_ = exec.Command(cliOpener, "http://localhost"+addr).Run()
	// done, serve types
	for {
		// cycle through these values...
		typechan <- &reflect.StructField{}
		typechan <- complex128(0)
		typechan <- reflect.Value{}
		typechan <- ""
	}
}
