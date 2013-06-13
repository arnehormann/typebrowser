package main

import (
	"github.com/arnehormann/typebrowser"
	"os/exec"
	"reflect"
	"runtime"
)

// Overdoing it MythBusters style...
type embedded0 struct{}
type embedded1 struct {
	a uint8
	b uint16
	c uint32
	d uint64
	e uint
	f func(interface {
		Func(uint)
	})
}
type embedded2 struct {
	a int8
	b int16
	c int32
	d int64
	e int `etag`
}
type compoundTest struct {
	embedded0
	embedded1
	a error `atag`
	_ [][][2]byte
	b *map[rune]*<-chan [2]uintptr
	c struct {
		a *complex64
		b complex128
		c interface{} `ctag`
		d interface {
			Do1()
			Do2() uintptr
			Do3(func() error)
			Do4(<-chan [2]uintptr, chan<- [2]uintptr) (bool, float32, float64)
		}
		e func(string, int) (bool, uint16)
		f map[struct{}]interface{}
	}
	embedded2
	d struct{}
}

func main() {
	addr := ":8080"
	typechan := typebrowser.NewTypeServer(addr)
	// open in browser
	cliOpener := "open"
	if runtime.GOOS == "windows" {
		cliOpener = "start"
	}
	_ = exec.Command(cliOpener, "http://localhost"+addr).Run()
	// done, serve types
	for {
		// cycle through these values...
		typechan <- typebrowser.Type{0, `example.go:64 - should be int`}
		typechan <- &reflect.StructField{}
		typechan <- complex128(0)
		typechan <- reflect.Value{}
		typechan <- typebrowser.Type{compoundTest{}, `example.go:68 - who shot my browser?`}
	}
}
