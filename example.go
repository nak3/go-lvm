package main

import (
	"fmt"
	//	"time"
)

func main() {
	vglist := ListVgNames()
	ListVgUUIDs()
	LvmPvListGet()
	a := &vgObject{}
	//a.vgt = LvmVgOpen(vglist[1], "r")
	a.vgt = VgOpen(vglist[1], "w")

	fmt.Printf("size: %d GiB\n", uint64(a.GetSize())/1024/1024/1024)

	fmt.Printf("pvlist: %#v\n", a.PvList())

	fmt.Printf("listLVs: %#v\n", a.LvList())

	// TODO /1024
	fmt.Printf("Free size: %d KiB\n", uint64(a.GetFreeSize())/1024/1024)

	l := &lvObject{}

	l = a.CreateLvLinear("foo", int64(a.GetFreeSize())/1024/1024/2)

	fmt.Printf("LV UUID: %#v\n", l.getUuid())

	l.addTag("Demo_tag")

	//	time.Sleep(10 * time.Second) // 3秒休む
	l.removeTag("Demo_tag")

	l.RemoveLv()

}
