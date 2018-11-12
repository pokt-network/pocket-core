package util

import "runtime"

// MyCaller returns the caller of the function that called it :)
func MyCaller() (*runtime.Func,uintptr){

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return nil,fpcs[0] // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0])
	if fun == nil {
		return nil,fpcs[0]
	}

	// return its name
	return fun, fpcs[0]
}
