package funcmock

import (
	"reflect"
	"sync"
)

type call struct {
	paramMutex *sync.Mutex
	// what parameter the function was passed
	param []interface{}

	// what the call to the function returned
	yield []interface{}

	// true if the call has been called, else false
	called bool
}

// func (this *call) AllParam() interface{} {
// 	return this.getParams()

// }

func (this *call) NthParam(nth int) interface{} {
	return this.getParams()[nth]
}

func (this *call) Called() bool {
	return this.called
}

func (this *call) setCalled(called bool) *call {
	this.called = called
	return this
}

func (this *call) SetReturn(args ...interface{}) *call {

	// Clears any previous return that was set

	this.yield = this.yield[:0]
	for _, nthIndex := range args {
		this.yield = append(this.yield, nthIndex)
	}
	return this
}

func (this *call) Return(args ...interface{}) *call {

	// needs to be thought over. Whether this function is required or not

	return this
}

func (this *call) getParams() (param []interface{}) {
	if this.Called() {
		this.paramMutex.Lock()
		param = this.param
		this.paramMutex.Unlock()
		return
	} else {
		panic("The nth call to the mock function has not been made yet")
	}
}

func (this *call) updateParam(param []interface{}) {
	this.setCalled(true)
	this.paramMutex.Lock()
	this.param = param
	this.paramMutex.Unlock()
	return
}

func sanitizeReturn(returnType reflect.Type, yield interface{}) (sanitizedYield reflect.Value) {

	if yield == nil {
		// kind of return param, eg. ptr, slice, etc.
		kind := returnType.Kind()
		switch kind {
		case reflect.Ptr:
		case reflect.Chan:
		case reflect.Func:
		case reflect.Interface:
		case reflect.Map:
		case reflect.Slice:
		default:
			panic("Cannot set nil to not-pointer type")
		}
		v := reflect.Zero(returnType)
		return v
	} else {

		return reflect.ValueOf(yield).Convert(returnType)
	}

}
