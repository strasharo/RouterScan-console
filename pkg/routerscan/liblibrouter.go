package routerscan

/*
#cgo LDFLAGS: -L${SRCDIR} -llibrouter
#include <stdlib.h>
#include <liblibrouter.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

func Initialize() error {
	if C.Initialize() == 0 {
		return errors.New("not initialized")
	}
	return nil
}

func GetModuleCount() (int, error) {
	var count C.dword
	if C.GetModuleCount(&count) == 0 {
		return 0, errors.New("cannot get module count")
	}
	return int(count), nil
}

type TModuleDesc struct {
	Enabled bool
	Name    string
	Desc    string
}

func charsToString(source []C.char) string {
	buf := make([]byte, 0)
	for i := 0; i < len(source); i++ {
		if byte(source[i]) != '\x00' {
			buf = append(buf, byte(source[i]))
		}
	}
	return string(buf)
}

func GetModuleInfo(index int) (*TModuleDesc, error) {
	desc := C.t_module_desc{}
	if C.GetModuleInfoW(C.uint(index), &desc) == 0 {
		return nil, errors.New("cannot get module info")
	}
	return &TModuleDesc{
		Enabled: desc.enabled == -1,
		Name:    charsToString(desc.name[:]),
		Desc:    charsToString(desc.desc[:]),
	}, nil
}

func SwitchModule(index int, enabled bool) error {
	var cEnabled int = 0
	if enabled {
		cEnabled = -1
	}
	if C.SwitchModule(C.uint(index), C.int(cEnabled)) == 0 {
		return errors.New("cannot switch module")
	}
	return nil
}

type TableDataCallbackT func(row uint, name string, value string)

var tdCallback TableDataCallbackT

//export tableDataCallback
func tableDataCallback(row C.uint, name *C.char, value *C.char) {
	tdCallback(uint(row), cCharToString(name), cCharToString(value))
}

type WriteLogCallbackT func(str string, verbosity int)

var wlCallback WriteLogCallbackT

//export writeLogCallback
func writeLogCallback(str *C.char, verbosity C.byte) {
	wlCallback(cCharToString(str), int(verbosity))
}

// cCharToString - converts delphi's strings to Go strings.
// Delphi's strings use nullbyte to specify end of char, not end of whole string,
// so we should manually find sequence \x00\x00 to determine end of string.
func cCharToString(val *C.char) string {
	testLen := 1
	for {
		testLen *= 2
		test := []byte(C.GoStringN(val, C.int(testLen)))
		for i := 0; i < len(test)-1; i++ {
			if test[i] == '\x00' && test[i+1] == '\x00' {
				return C.GoStringN(val, C.int(i))
			}
		}
	}
}

func stringToCChar(val string) *C.char {
	buf := make([]byte, 0)
	runes := []byte(val)
	for i := 0; i < len(runes); i++ {
		buf = append(buf, byte(runes[i]), '\x00')
	}
	buf = append(buf, '\x00')
	return C.CString(string(buf))
}

func SetParamBool(option StValueBool, value bool) error {
	var cValue uintptr = 0
	if value {
		cValue = 3
	}

	if C.SetParamW(C.uint(option), unsafe.Pointer(cValue)) == 0 {
		return fmt.Errorf("cannot set param %d", option)
	}
	return nil
}

func SetParamString(option StValueString, value string) error {
	cValue := stringToCChar(value)
	//cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if C.SetParamW(C.uint(option), unsafe.Pointer(cValue)) == 0 {
		return fmt.Errorf("cannot set param %d", option)
	}
	return nil
}

func SetSetTableDataCallback(cb TableDataCallbackT) error {
	tdCallback = cb
	if C.SetParamW(C.uint(StSetTableDataCallback), unsafe.Pointer(C.tableDataCallback)) == 0 {
		return errors.New("cannot set TableDataCallback")
	}
	return nil
}

func SetWriteLogCallback(cb WriteLogCallbackT) error {
	wlCallback = cb
	if C.SetParamW(C.uint(StWriteLogCallback), unsafe.Pointer(C.writeLogCallback)) == 0 {
		return errors.New("cannot set WriteLogCallback")
	}
	return nil
}

type Router struct {
	p unsafe.Pointer
}

func PrepareRouter(row int, ip uint32, port uint16) (*Router, error) {
	var addr C.uint
	address := C.malloc(C.size_t(unsafe.Sizeof(addr)))
	if C.PrepareRouter(C.uint(row), C.uint(ip), C.ushort(port), unsafe.Pointer(&address)) == 0 {
		return nil, errors.New("cannot prepare router")
	}
	return &Router{
		p: address,
	}, nil
}

func (router Router) Scan() error {
	if C.ScanRouter(unsafe.Pointer(router.p)) == 0 {
		return errors.New("cannot run scan")
	}
	return nil
}

func (router Router) Free() error {
	if C.FreeRouter(unsafe.Pointer(router.p)) == 0 {
		return errors.New("cannot free router")
	}
	return nil
}
