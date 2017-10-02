package native

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/pkg/errors"
)

func Cfopen(filename string, mode string) (CFILE, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cMode := C.CString(mode)
	defer C.free(unsafe.Pointer(cMode))

	file, err := C.fopen(cFilename, cMode)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(file), nil
}

func Cfclose(file CFILE) error {
	res := C.fclose((*C.FILE)(file))
	if res != 0 {
		return errors.New("failed to close file")
	}
	return nil
}