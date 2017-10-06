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

	fmt.Printf("(test)pvsList: %#v\n", gs)
	return gs
}

func VgOpen(vgname string, mode string) C.vg_t {
	if mode == "" {
		mode = "r"
	}
	vg := C.lvm_vg_open(libh, C.CString(vgname), C.CString(mode), 0)
	return vg
}

type GoObject struct {
	vgt C.vg_t
	lvt C.lv_t
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

// TODO: should be for lvt
func (g *GoObject) getUuid() *C.char {
	return C.lvm_lv_get_uuid(g.lvt)
}

// VG methods
type vgObject struct {
	vgt C.vg_t
}

//getName
func (v *vgObject) getName() *C.char {
	return C.lvm_vg_get_name(v.vgt)
}

//getUuid
func (v *vgObject) getUuid() *C.char {
	return C.lvm_vg_get_uuid(v.vgt)
}

// close

// PvList lists of pvs from VG
func (v *vgObject) PvList() []string {
	pvs := C.lvm_vg_list_pvs(v.vgt)
	if pvs == nil {
		return []string{""}
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(pvs, cargs)
	gs := GoStrings(n, cargs)
	fmt.Printf("(test)pvsList: %#v\n", gs)
	return gs

}

// GetSize returns size of VG
func (v *vgObject) GetSize() C.uint64_t {
	return C.lvm_vg_get_size(v.vgt)
}

// LvList lists of lvs from VG
func (v *vgObject) LvList() []string {
	lvl := C.lvm_vg_list_lvs(v.vgt)
	if lvl == nil {
		return []string{""}
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(lvl, cargs)
	gs := GoStrings(n, cargs)
	fmt.Printf("(test)lvList: %#v\n", gs)
	return gs
}

// GetFreeSize returns free size of VG
func (v *vgObject) GetFreeSize() C.uint64_t {
	return C.lvm_vg_get_free_size(v.vgt)
}

func createGoLv(v *vgObject, lv C.lv_t) *lvObject {
	return &lvObject{
		lvt:      lv,
		parentVG: v,
	}
}

func (v *vgObject) CreateLvLinear(n string, s int64) *lvObject {
	//func (v *vgObject) CreateLvLinear(n string, s int64) *C.struct_logical_volume {
	size := C.uint64_t(s)
	name := C.CString(n)

	lv := C.lvm_vg_create_lv_linear(v.vgt, name, size)
	if lv == nil {
		fmt.Printf("nil")
	}
	return createGoLv(v, lv)
	//	return &lvObject{lvt: lv}
	//	return lv
}

// ######################################## LV ###################################

// LV object
type lvObject struct {
	lvt      C.lv_t
	parentVG *vgObject
}

// LV methods
// getAttr
func (l *lvObject) getAttr() *C.char {
	return C.lvm_lv_get_attr(l.lvt)
}

// getOrigin
func (l *lvObject) getOrigin() *C.char {
	return C.lvm_lv_get_origin(l.lvt)
}

// getName
func (l *lvObject) getName() *C.char {
	return C.lvm_lv_get_name(l.lvt)
}

// getUuid
func (l *lvObject) getUuid() *C.char {
	return C.lvm_lv_get_uuid(l.lvt)
}

// RemoveLV
func (l *lvObject) RemoveLv() error {
	// TODO return
	C.lvm_vg_remove_lv(l.lvt)
	return nil
}

// addTag
func (l *lvObject) addTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_add_tag(l.lvt, tag)
	C.lvm_vg_write(l.parentVG.vgt)
	return nil
}

// removeTag
func (l *lvObject) removeTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_remove_tag(l.lvt, tag)
	C.lvm_vg_write(l.parentVG.vgt)
	return nil
}

// ######################## pv methods #######################

// pvObject
type pvObject struct {
	pvt C.pv_t
}

// getName
func (p *pvObject) getName() *C.char {
	return C.lvm_pv_get_name(p.pvt)
}

// getUuid
func (p *pvObject) getUuid() *C.char {
	return C.lvm_pv_get_uuid(p.pvt)
}

// getMdaCount
func (p *pvObject) getMdaCount() C.uint64_t {
	return C.lvm_pv_get_mda_count(p.pvt)
}
