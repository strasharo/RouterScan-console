package routerscan

import "C"

func charsToString(source []C.char) string {
	buf := make([]byte, 0)
	for i := 0; i < len(source); i++ {
		if byte(source[i]) != '\x00' {
			buf = append(buf, byte(source[i]))
		}
	}
	return string(buf)
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
