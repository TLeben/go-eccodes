package native

/*
#include <stdio.h>
*/
import "C"

import (
	"unsafe"

	"github.com/tleben/go-errors/v2"
)

func Cfopen(filename string, mode string) (CFILE, error) {
	cFilename := C.CString(filename)
	defer Cfree(unsafe.Pointer(cFilename))

	cMode := C.CString(mode)
	defer Cfree(unsafe.Pointer(cMode))

	file, err := C.fopen(cFilename, cMode)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(file), nil
}

func Cfclose(file CFILE) error {
	res := C.fclose((*C.FILE)(file))
	if res != 0 {
		return errors.Error("failed to close io")
	}
	return nil
}

func Cmclose(memory Cbytes) error {
	Cfree(memory)
	return nil
}
