package native

/*
#include <eccodes.h>
*/
import "C"

import (
	"io"
	"unsafe"

	"github.com/tleben/go-errors/v2"
)

func Ccodes_handle_new_from_index(index Ccodes_index) (Ccodes_handle, error) {
	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	h := C.codes_handle_new_from_index((*C.codes_index)(index), cError)
	if err != 0 {
		if err == Cint(C.CODES_END_OF_INDEX) {
			return nil, io.EOF
		}
		return nil, errors.Error(Cgrib_get_error_message(int(err)))
	}
	return unsafe.Pointer(h), nil
}

func Ccodes_handle_new_from_file(ctx Ccodes_context, file CFILE, product int) (Ccodes_handle, error) {
	var cProduct C.int

	cProduct = C.int(product)

	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	h := C.codes_handle_new_from_file((*C.grib_context)(ctx), (*C.FILE)(file), C.ProductKind(cProduct), cError)
	if err != 0 {
		return nil, errors.Error(Cgrib_get_error_message(int(err)))
	}

	if h == nil {
		return nil, io.EOF
	}

	return unsafe.Pointer(h), nil
}

func Ccodes_handle_new_from_message(ctx Ccodes_context, data Cbytes, size Clong) (Ccodes_handle, error) {
	h := C.codes_handle_new_from_message((*C.grib_context)(ctx), data, C.ulong(size))

	//	Csize := C.ulong(size)

	//	var err Cint
	//	cError := (*C.int)(unsafe.Pointer(&err))
	//	h := C.codes_grib_handle_new_from_multi_message((*C.grib_context)(ctx), &data, &Csize, cError)

	//	if err != 0 {
	//		return nil, errors.Error(Cgrib_get_error_message(int(err)))
	//	}

	if h == nil {
		return nil, io.EOF
	}

	return unsafe.Pointer(h), nil
}

func Ccodes_handle_delete(handle Ccodes_handle) error {
	err := C.codes_handle_delete((*C.codes_handle)(handle))
	if err != 0 {
		return errors.Error(Cgrib_get_error_message(int(err)))
	}
	return nil
}
