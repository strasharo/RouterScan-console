package routerscan

/*
#cgo LDFLAGS: -L${SRCDIR} -llibrouter
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

type ParamBool int

const (
	StEnableDebug ParamBool = 0
)

func SetParamBool(st ParamBool, value bool) error {
	var cValue C.uint = 1
	if value {
		cValue = 0
	}
	if C.SetParamW(C.uint(st), unsafe.Pointer(&cValue)) == 0 {
		return fmt.Errorf("cannot set bool param %d", st)
	}
	return nil
}

//export tableDataCallback
func tableDataCallback(row C.uint, name *C.char, value *C.char) {
	fmt.Println(uint(row), cCharToString(name), cCharToString(value))
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

func SetSetTableDataCallback() error {
	if C.SetParamW(C.uint(3), unsafe.Pointer(C.tableDataCallback)) == 0 {
		return errors.New("cannot set TableDataCallback")
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
