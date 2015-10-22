package funcmock

import "reflect"

type call struct {

	// what parameter the function was passed
	param []interface{}

	// what the call to the function returned
	yield []interface{}

	// true if the call has been called, else false
	called bool
}

func (this *call) ParamNth(nth int) interface{} {
	return this.param[nth]
}

func (this *call) Called() bool {
	return this.called
}

func (this *call) setCalled(called bool) *call {
	this.called = true
	return this
}

func (this *call) SetReturn(args ...interface{}) *call {

	/*

		Clears any previous return that was set

	*/
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
