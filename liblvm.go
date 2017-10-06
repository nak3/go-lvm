package lvm

//#cgo LDFLAGS: -llvm2app
//#include "macro_wrapper.h"
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

func GoStrings(argc C.int, argv **C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

// VG methods
type VgObject struct {
	Vgt C.vg_t
}

//getName
func (v *VgObject) GetName() *C.char {
	return C.lvm_vg_get_name(v.Vgt)
}

//getUuid
func (v *VgObject) GetUuid() *C.char {
	return C.lvm_vg_get_uuid(v.Vgt)
}

// close

// PvList lists of pvs from VG
func (v *VgObject) PvList() []string {
	pvs := C.lvm_vg_list_pvs(v.Vgt)
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
func (v *VgObject) GetSize() C.uint64_t {
	return C.lvm_vg_get_size(v.Vgt)
}

// LvList lists of lvs from VG
func (v *VgObject) LvList() []string {
	lvl := C.lvm_vg_list_lvs(v.Vgt)
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
func (v *VgObject) GetFreeSize() C.uint64_t {
	return C.lvm_vg_get_free_size(v.Vgt)
}

func createGoLv(v *VgObject, lv C.lv_t) *LvObject {
	return &LvObject{
		Lvt:      lv,
		parentVG: v,
	}
}

func (v *VgObject) CreateLvLinear(n string, s int64) *LvObject {
	//func (v *VgObject) CreateLvLinear(n string, s int64) *C.struct_logical_volume {
	size := C.uint64_t(s)
	name := C.CString(n)

	lv := C.lvm_vg_create_lv_linear(v.Vgt, name, size)
	if lv == nil {
		fmt.Printf("nil")
	}
	return createGoLv(v, lv)
	//	return &LvObject{Lvt: lv}
	//	return lv
}

// ######################################## LV ###################################

// LV object
type LvObject struct {
	Lvt      C.lv_t
	parentVG *VgObject
}

// LV methods
// getAttr
func (l *LvObject) getAttr() *C.char {
	return C.lvm_lv_get_attr(l.Lvt)
}

// getOrigin
func (l *LvObject) getOrigin() *C.char {
	return C.lvm_lv_get_origin(l.Lvt)
}

// getName
func (l *LvObject) getName() *C.char {
	return C.lvm_lv_get_name(l.Lvt)
}

// getUuid
func (l *LvObject) GetUuid() *C.char {
	return C.lvm_lv_get_uuid(l.Lvt)
}

// RemoveLV
func (l *LvObject) RemoveLv() error {
	// TODO return
	C.lvm_vg_remove_lv(l.Lvt)
	return nil
}

// addTag
func (l *LvObject) AddTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_add_tag(l.Lvt, tag)
	C.lvm_vg_write(l.parentVG.Vgt)
	return nil
}

// removeTag
func (l *LvObject) RemoveTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_remove_tag(l.Lvt, tag)
	C.lvm_vg_write(l.parentVG.Vgt)
	return nil
}

// ######################## pv methods #######################

// pvObject
type pvObject struct {
	pvt C.pv_t
}

// getName
func (p *pvObject) GetName() *C.char {
	return C.lvm_pv_get_name(p.pvt)
}

// getUuid
func (p *pvObject) GetUuid() *C.char {
	return C.lvm_pv_get_uuid(p.pvt)
}

// getMdaCount
func (p *pvObject) getMdaCount() C.uint64_t {
	return C.lvm_pv_get_mda_count(p.pvt)
}
