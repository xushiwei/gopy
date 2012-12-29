// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void xdecref(PyObject *obj) { Py_XDECREF(obj); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Error struct {
	kind *C.PyObject
	tb   *C.PyObject
	msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("<TYPE>: %s", e.msg)
}

func exceptionRaised() bool {
	return C.PyErr_Occurred() != nil
}

func exception() error {
	if C.PyErr_Occurred() == nil {
		return nil
	}

	var t, v, tb *C.PyObject

	defer C.xdecref(v)

	C.PyErr_Fetch(&t, &v, &tb)

	s := C.PyObject_Str(v)
	defer C.decref(s)

	return &Error{t, tb, C.GoString(C.PyString_AsString(s))}
}

func raise(err error) {
	var msg string
	var exc = C.PyExc_Exception

	e, ok := err.(*Error)
	if ok {
		exc = e.kind
		msg = e.msg
	} else {
		msg = err.Error()
	}

	s := C.CString(msg)
	defer C.free(unsafe.Pointer(s))

	C.PyErr_SetString(exc, s)
}

func Err_Format(f string, args ...interface{}) {
	s := C.CString(fmt.Sprintf(f, args))
	defer C.free(unsafe.Pointer(s))
	C.PyErr_SetString(C.PyExc_Exception, s)
}

func TypeError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	C.incref(C.PyExc_TypeError)
	return &Error{C.PyExc_TypeError, nil, msg}
}

func KeyError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	C.incref(C.PyExc_KeyError)
	return &Error{C.PyExc_KeyError, nil, msg}
}

func AttributeError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	C.incref(C.PyExc_AttributeError)
	return &Error{C.PyExc_AttributeError, nil, msg}
}

func NotImplemented(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	C.incref(C.PyExc_NotImplementedError)
	return &Error{C.PyExc_NotImplementedError, nil, msg}
}

func int2Err(i C.int) error {
	if i < 0 {
		return exception()
	}
	return nil
}

func int2BoolErr(i C.int) (bool, error) {
	if i < 0 {
		return false, exception()
	}
	return i > 0, nil
}

func ssize_t2Int64Err(s C.Py_ssize_t) (int64, error) {
	if s < 0 {
		return 0, exception()
	}
	return int64(s), nil
}

func obj2ObjErr(obj *C.PyObject) (Object, error) {
	if obj == nil {
		return nil, exception()
	}
	return newObject(obj), nil
}
