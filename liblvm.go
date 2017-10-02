package main

//#cgo LDFLAGS: -llvm2app
//#include "macro_wrapper.c"
import "C"
import (
	"fmt"
	"unsafe"
)

var libh *C.struct_lvm

func init() {
	libh = C.lvm_init(nil)
}

// ListVgNames returns LVM name list
func ListVgNames() []string {
	vgnames := C.lvm_list_vg_names(libh)
	if vgnames != nil {
		// TODO
		fmt.Printf("nil\n")
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(vgnames, cargs)
	gs := GoStrings(n, cargs)
	//
	fmt.Printf("vgnames: %#v\n", gs)
	return gs
}

// ListVgUUIDs returns LVM uuid list
func ListVgUUIDs() []string {
	uuids := C.lvm_list_vg_uuids(libh)
	if uuids != nil {
		// TODO
		fmt.Printf("nil\n")
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(uuids, cargs)
	gs := GoStrings(n, cargs)
	//
	fmt.Printf("vgnames: %#v\n", gs)
	return gs
}

func LvmPvListGet() []string {
	pvsList := C.lvm_list_pvs(libh)
	fmt.Printf("pvsList: %#v\n", pvsList)

	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(pvsList, cargs)
	gs := GoStrings(n, cargs)

	fmt.Printf("pvsList: %#v\n", gs)
	return gs
}

func LvmVgOpen(vgname string, mode string) C.vg_t {
	if mode == "" {
		mode = "r"
	}
	vg := C.lvm_vg_open(libh, C.CString(vgname), C.CString(mode), 0)
	return vg

}

type GoObject struct {
}

// uint 64 for size?
func CreateLvLinear(vgobject *C.struct_volume_group, n string, s int64) GoObject {
	size := C.uint64_t(s)
	name := C.CString(n)

	lv := C.lvm_vg_create_lv_linear(vgobject, name, size)
	//	vg := LvmVgOpen("foo", "r")
	//	a := []string{"fo"}
	return create_go_lv(vgobject, lv)
}

//type C.vg_t C.struct_volume_group

func create_go_lv(vgobject *C.struct_volume_group, lv C.lv_t) GoObject {
	return GoObject{
	//	TODO
	}
}

func GoStrings(argc C.int, argv **C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}
