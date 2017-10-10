package lvm_test

import (
	"fmt"

	"github.com/nak3/go-lvm"
)

// This example demonstrates creating and removing a Logical Volume(LV).
func ExampleLvObject_createremove() {
	// List volume group
	vglist := lvm.ListVgNames()
	availableVG := ""

	// Create a VG object
	vgo := &lvm.VgObject{}
	for i := 0; i < len(vglist); i++ {
		vgo.Vgt = lvm.VgOpen(vglist[i], "r")
		if vgo.GetFreeSize() > 0 {
			availableVG = vglist[i]
			vgo.Close()
			break
		}
		vgo.Close()
	}
	if availableVG == "" {
		fmt.Printf("no VG that has free space found\n")
		return
	}

	// Open VG in write mode
	vgo.Vgt = lvm.VgOpen(availableVG, "w")
	defer vgo.Close()

	// Create a LV object
	l := &lvm.LvObject{}

	// Create a LV
	l, err := vgo.CreateLvLinear("go-lvm-example-test-lv", int64(vgo.GetFreeSize())/1024/1024/2)
	if err != nil {
		fmt.Printf("error: %v")
		return
	}

	// Output uuid of LV
	fmt.Printf("Created\n\tuuid: %s\n\tname: %s\n\tattr: %s\n\torigin: %s\n",
		l.GetUuid(), l.GetName(), l.GetAttr(), l.GetOrigin())
	// Output uuid of LV
	l.Remove()

	/*
	   Created
	   	uuid: cn631J-J2GR-DL0l-3G38-MfGm-8ypc-iHskGI
	   	name: go-lvm-example-test-lv
	   	attr: -wi-a-----
	   	origin:
	*/
}
