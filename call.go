package funcmock

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

func (this *call) SetCalled(called bool) *call {
	this.called = true
	return this
}

func (this *call) SetReturn(args ...interface{}) *call {

	/*

		Clears any previous return that was set

	*/

	for _, nthIndex := range args {
		this.yield = append(this.yield, nthIndex)
	}
	return this
}

func (this *call) Return(args ...interface{}) *call {

	// needs to be thought over. Whether this function is required or not

	return this
}
