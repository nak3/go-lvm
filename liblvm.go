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

// ######################## LVM methods #######################

// GetVersion returns library version
func GetVersion() *C.char {
	return C.lvm_library_get_version()
}

// GC cleans up libh
func GC() {
	C.lvm_quit(libh)
	libh = nil
}

//VgOpen opens volume group
func VgOpen(vgname string, mode string) C.vg_t {
	if mode == "" {
		mode = "r"
	}
	vg := C.lvm_vg_open(libh, C.CString(vgname), C.CString(mode), 0)
	return vg
}

//VgCreate creates VG
func VgCreate(vgname string) C.vg_t {
	return C.lvm_vg_create(libh, C.CString(vgname))
}

//        { "configFindBool",     (PyCFunction)_liblvm_lvm_config_find_bool, METH_VARARGS },
//        { "configReload",       (PyCFunction)_liblvm_lvm_config_reload, METH_NOARGS },
//        { "configOverride",     (PyCFunction)_liblvm_lvm_config_override, METH_VARARGS },

// Scan scans libh
func Scan() {
	C.lvm_scan(libh)
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
	gs := goStrings(n, cargs)
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
	gs := goStrings(n, cargs)
	return gs
}

//PvCreate creates PV
func PvCreate(pv_name string, size uint64, pvmetadatacopies uint64, pvmetadatasize uint64,
	data_alignment uint64, data_alignment_offset uint64, zero uint64) C.vg_t {
	pv_params := C.lvm_pv_params_create(libh, C.CString(pv_name))
	if pv_params != nil {
		// TODO
	}

	// TODO return error
	C.wrapper_set_pv_prop(pv_params, C.CString("size"), C.longlong(size))
	C.wrapper_set_pv_prop(pv_params, C.CString("size"), C.longlong(size))
	C.wrapper_set_pv_prop(pv_params, C.CString("pvmetadatacopies"), C.longlong(pvmetadatacopies))
	C.wrapper_set_pv_prop(pv_params, C.CString("pvmetadatasize"), C.longlong(pvmetadatasize))
	C.wrapper_set_pv_prop(pv_params, C.CString("data_alignment"), C.longlong(data_alignment))
	C.wrapper_set_pv_prop(pv_params, C.CString("data_alignment_offset"), C.longlong(data_alignment_offset))
	C.wrapper_set_pv_prop(pv_params, C.CString("zero"), C.longlong(zero))

	C.lvm_pv_create_adv(pv_params)
	// TODO
	return nil
}

// PvRemove removes PV
func PvRemove(pv_name string) {
	C.lvm_pv_remove(libh, C.CString(pv_name))
}

// PercentToFloat converts percent to float
func PercentToFloat(percent C.percent_t) C.float {
	// TODO C.percent_t should be golang type.
	return C.lvm_percent_to_float(percent)
}

// VgNameValidate validates if the vg name is valid
func VgNameValidate(name string) error {
	ret := C.lvm_vg_name_validate(libh, C.CString(name))
	if ret < 0 {
		return getLastError()
	}
	return nil
}

// VgNameFromPvID returns VG name from PVID
func VgNameFromPvID(pvid string) *C.char {
	ret := C.lvm_vgname_from_pvid(libh, C.CString(pvid))
	// TODO
	msg := C.lvm_errmsg(libh)
	fmt.Printf("msg : %#v\n", C.GoString(msg))
	return ret
}

// VgNameFromPvDevice returns VG name from PV device
func VgNameFromPvDevice(pvdevice string) string {
	//func VgNameFromPvDevice(pvdevice string) *C.char {
	ret := C.lvm_vgname_from_device(libh, C.CString(pvdevice))
	// TODO
	msg := C.lvm_errmsg(libh)
	fmt.Printf("msg : %#v\n", C.GoString(msg))
	return C.GoString(ret)
}

// ######################## vg methods #######################

// VgObject is an object of VG
type VgObject struct {
	Vgt C.vg_t
}

// GetName gets name of VG
func (v *VgObject) GetName() string {
	return C.GoString(C.lvm_vg_get_name(v.Vgt))
}

// GetUuid gets UUID of VG
func (v *VgObject) GetUuid() *C.char {
	return C.lvm_vg_get_uuid(v.Vgt)
}

// Close closes VG object
func (v *VgObject) Close() error {
	if C.lvm_vg_close(v.Vgt) == -1 {
		return getLastError()
	}
	return nil
}

// Remove removes VG
func (v *VgObject) Remove() error {
	if C.lvm_vg_remove(v.Vgt) == -1 {
		return getLastError()
	}
	if C.lvm_vg_write(v.Vgt) == -1 {
		return getLastError()
	}
	return v.Close()
}

//        { "extend",             (PyCFunction)_liblvm_lvm_vg_extend, METH_VARARGS },
//        { "reduce",             (PyCFunction)_liblvm_lvm_vg_reduce, METH_VARARGS },
//        { "addTag",             (PyCFunction)_liblvm_lvm_vg_add_tag, METH_VARARGS },
//        { "removeTag",          (PyCFunction)_liblvm_lvm_vg_remove_tag, METH_VARARGS },
//        { "setExtentSize",      (PyCFunction)_liblvm_lvm_vg_set_extent_size, METH_VARARGS },
//        { "isClustered",        (PyCFunction)_liblvm_lvm_vg_is_clustered, METH_NOARGS },
//        { "isExported",         (PyCFunction)_liblvm_lvm_vg_is_exported, METH_NOARGS },
//        { "isPartial",          (PyCFunction)_liblvm_lvm_vg_is_partial, METH_NOARGS },
//        { "getSeqno",           (PyCFunction)_liblvm_lvm_vg_get_seqno, METH_NOARGS },

// GetSize returns size of VG
func (v *VgObject) GetSize() C.uint64_t {
	return C.lvm_vg_get_size(v.Vgt)
}

// GetFreeSize returns free size of VG
func (v *VgObject) GetFreeSize() C.uint64_t {
	return C.lvm_vg_get_free_size(v.Vgt)
}

//        { "getExtentSize",      (PyCFunction)_liblvm_lvm_vg_get_extent_size, METH_NOARGS },
//        { "getExtentCount",     (PyCFunction)_liblvm_lvm_vg_get_extent_count, METH_NOARGS },
//        { "getFreeExtentCount", (PyCFunction)_liblvm_lvm_vg_get_free_extent_count, METH_NOARGS },
//        { "getProperty",        (PyCFunction)_liblvm_lvm_vg_get_property, METH_VARARGS },
//        { "setProperty",        (PyCFunction)_liblvm_lvm_vg_set_property, METH_VARARGS },
//        { "getPvCount",         (PyCFunction)_liblvm_lvm_vg_get_pv_count, METH_NOARGS },
//        { "getMaxPv",           (PyCFunction)_liblvm_lvm_vg_get_max_pv, METH_NOARGS },
//        { "getMaxLv",           (PyCFunction)_liblvm_lvm_vg_get_max_lv, METH_NOARGS },

// ListLVs lists of lvs from VG
func (v *VgObject) ListLVs() []string {
	lvl := C.lvm_vg_list_lvs(v.Vgt)
	if lvl == nil {
		return []string{""}
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(lvl, cargs)
	gs := goStrings(n, cargs)
	fmt.Printf("(test)lvList: %#v\n", gs)
	return gs
}

// ListPVs lists of pvs from VG
func (v *VgObject) ListPVs() []string {
	pvs := C.lvm_vg_list_pvs(v.Vgt)
	if pvs == nil {
		return []string{""}
	}
	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(pvs, cargs)
	gs := goStrings(n, cargs)
	fmt.Printf("(test)pvsList: %#v\n", gs)
	return gs
}

//        { "lvFromName",         (PyCFunction)_liblvm_lvm_lv_from_name, METH_VARARGS },
//        { "lvFromUuid",         (PyCFunction)_liblvm_lvm_lv_from_uuid, METH_VARARGS },
//        { "lvNameValidate",     (PyCFunction)_liblvm_lvm_lv_name_validate, METH_VARARGS },
//        { "pvFromName",         (PyCFunction)_liblvm_lvm_pv_from_name, METH_VARARGS },
//        { "pvFromUuid",         (PyCFunction)_liblvm_lvm_pv_from_uuid, METH_VARARGS },
//        { "getTags",            (PyCFunction)_liblvm_lvm_vg_get_tags, METH_NOARGS },

// createGoLv creats a LV Object
func createGoLv(v *VgObject, lv C.lv_t) *LvObject {
	return &LvObject{
		Lvt:      lv,
		parentVG: v,
	}
}

// CreateLvLinear creates LV Object
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

//        { "createLvThinpool",   (PyCFunction)_liblvm_lvm_vg_create_lv_thinpool, METH_VARARGS },
//        { "createLvThin",       (PyCFunction)_liblvm_lvm_vg_create_lv_thin, METH_VARARGS },

// ######################################## LV methods ###################################

// LV object
type LvObject struct {
	Lvt      C.lv_t
	parentVG *VgObject
}

// GetAttr gets LV attr
func (l *LvObject) GetAttr() string {
	return C.GoString(C.lvm_lv_get_attr(l.Lvt))
}

// GetName gets LV name
func (l *LvObject) GetName() string {
	return C.GoString(C.lvm_lv_get_name(l.Lvt))
}

// GetOrigin gets LV origin
func (l *LvObject) GetOrigin() string {
	return C.GoString(C.lvm_lv_get_origin(l.Lvt))
}

// GetUuid gets LV uuid
func (l *LvObject) GetUuid() string {
	return C.GoString(C.lvm_lv_get_uuid(l.Lvt))
}

//        { "activate",           (PyCFunction)_liblvm_lvm_lv_activate, METH_NOARGS },
//        { "deactivate",         (PyCFunction)_liblvm_lvm_lv_deactivate, METH_NOARGS },

// Remove removes LV
func (l *LvObject) Remove() error {
	// TODO return
	C.lvm_vg_remove_lv(l.Lvt)
	return nil
}

//        { "getProperty",        (PyCFunction)_liblvm_lvm_lv_get_property, METH_VARARGS },
//        { "getSize",            (PyCFunction)_liblvm_lvm_lv_get_size, METH_NOARGS },
//        { "isActive",           (PyCFunction)_liblvm_lvm_lv_is_active, METH_NOARGS },
//        { "isSuspended",        (PyCFunction)_liblvm_lvm_lv_is_suspended, METH_NOARGS },

// AddTag adds tag to LV
func (l *LvObject) AddTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_add_tag(l.Lvt, tag)
	C.lvm_vg_write(l.parentVG.Vgt)
	return nil
}

// RemoveTag removes tag from LV
func (l *LvObject) RemoveTag(stag string) error {
	tag := C.CString(stag)
	C.lvm_lv_remove_tag(l.Lvt, tag)
	C.lvm_vg_write(l.parentVG.Vgt)
	return nil
}

//        { "getTags",            (PyCFunction)_liblvm_lvm_lv_get_tags, METH_NOARGS },
//        { "rename",             (PyCFunction)_liblvm_lvm_lv_rename, METH_VARARGS },
//        { "resize",             (PyCFunction)_liblvm_lvm_lv_resize, METH_VARARGS },
//        { "listLVsegs",         (PyCFunction)_liblvm_lvm_lv_list_lvsegs, METH_NOARGS },
//        { "snapshot",           (PyCFunction)_liblvm_lvm_lv_snapshot, METH_VARARGS },

// ######################## pv list methods #######################

// Open lists PVs and get them as string array
func Open() []string {
	pvsList := C.lvm_list_pvs(libh)
	fmt.Printf("pvsList: %#v\n", pvsList)

	cargs := C.makeCharArray(C.int(0))
	n := C.wrapper_dm_list_iterate_items(pvsList, cargs)
	gs := goStrings(n, cargs)
	return gs
}

//        { "close",              (PyCFunction)_liblvm_lvm_pvlist_put, METH_VARARGS },

// ######################## pv methods #######################

// pvObject
type pvObject struct {
	pvt C.pv_t
}

// getName
func (p *pvObject) GetName() string {
	return C.GoString(C.lvm_pv_get_name(p.pvt))
}

// getUuid
func (p *pvObject) GetUuid() string {
	return C.GoString(C.lvm_pv_get_uuid(p.pvt))
}

// getMdaCount
func (p *pvObject) getMdaCount() C.uint64_t {
	return C.lvm_pv_get_mda_count(p.pvt)
}

//        { "getProperty",        (PyCFunction)_liblvm_lvm_pv_get_property, METH_VARARGS },
//        { "getSize",            (PyCFunction)_liblvm_lvm_pv_get_size, METH_NOARGS },
//        { "getDevSize",         (PyCFunction)_liblvm_lvm_pv_get_dev_size, METH_NOARGS },
//        { "getFree",            (PyCFunction)_liblvm_lvm_pv_get_free, METH_NOARGS },
//        { "resize",             (PyCFunction)_liblvm_lvm_pv_resize, METH_VARARGS },
//        { "listPVsegs",         (PyCFunction)_liblvm_lvm_pv_list_pvsegs, METH_NOARGS },

// ######################## utility methods #######################
func goStrings(argc C.int, argv **C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

func getLastError() error {
	msg := C.GoString(C.lvm_errmsg(libh))
	if msg == "" {
		return fmt.Errorf("unknown error")
	}
	return fmt.Errorf(msg)
}
