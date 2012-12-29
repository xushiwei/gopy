// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

//export goClassCall
func goClassCall(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) (Object, error))(unsafe.Pointer(&ctxt.call))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	ret, err := (*f)(obj, a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCompare
func goClassCompare(obj1, obj2 unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (int, error))(unsafe.Pointer(&ctxt.compare))

	o := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, o)
	if err != nil {
		raise(err)
		return -1
	}

	return ret
}

//export goClassDealloc
func goClassDealloc(obj unsafe.Pointer) {
	// Get the class context
	ctxt := getClassContext(obj)

	if ctxt.dealloc != nil {
		// Turn the function into something we can call
		f := (*func(unsafe.Pointer))(unsafe.Pointer(&ctxt.dealloc))

		(*f)(obj)
	} else {
		(*AbstractObject)(obj).Free()
	}
}

//export goClassInit
func goClassInit(obj, args, kwds unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) error)(unsafe.Pointer(&ctxt.init))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	err := (*f)(obj, a, k)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassRepr
func goClassRepr(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.repr))

	s := C.CString((*f)(obj))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyString_FromString(s))
}

//export goClassRichCmp
func goClassRichCmp(obj1, obj2 unsafe.Pointer, op int) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object, Op) (Object, error))(unsafe.Pointer(&ctxt.richcmp))

	// Get obj2 ready for use
	arg := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, arg, Op(op))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassStr
func goClassStr(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.str))

	s := C.CString((*f)(obj))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyString_FromString(s))
}
