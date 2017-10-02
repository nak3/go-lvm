package main

import (
	"fmt"
)

func main() {
	ListVgNames()
	ListVgUUIDs()
	LvmPvListGet()
	a := LvmVgOpen("test", "r")
	fmt.Printf("%#v", a)
}
