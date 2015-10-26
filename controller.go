package funcmock

import "reflect"

type MockController struct {
	originalFunc reflect.Value
	targetFunc   reflect.Value
	// we need map, not slice, to set call before it is called
	callStack chan map[int]*call
	// we need it to set call, before it is called
	counter chan int
	// the default call which shall be used for mint calls
	defaultYield []reflect.Value

	// Flag indicating the default return has been set
	yieldSet bool
}

func (this *MockController) CallCount() int {
	counter := <-this.counter
	go func() { this.counter <- counter }()
	return counter
}

func (this *MockController) NthCall(nth int) (c *call) {
	c = this.callStack[nth]
	if c == nil {
		c = new(call)
		this.callStack[nth] = c
	}
	return c
}

func (this *MockController) incrementCounter() {
	counter := <-this.counter
	counter++
	go func() { this.counter <- counter }()
}

func (this *MockController) SetDefaultReturn(args ...interface{}) {
	if this.targetFunc == reflect.Zero(this.targetFunc.Type()) {
		panic("Internal Error: Target Function should prior to calling SetDefaultReturn")
	}
	fnType := this.targetFunc.Type()
	typeNumOut := fnType.NumOut()
	if len(args) == typeNumOut && !this.yieldSet {
		this.defaultYield = this.defaultYield[:0]
		for i := 0; i < typeNumOut; i++ {
			this.defaultYield = append(this.defaultYield, sanitizeReturn(fnType.Out(i), args[i]))
		}
		this.yieldSet = true
	} else if this.yieldSet {
		panic("Can only call SetDefaultReturn once")
	} else {
		panic("The number of returns should be the same as that of the function")
	}

}

func (this *MockController) add(c *call) {
	this.callStack[this.CallCount()-1] = c
}

func (this *MockController) Called() bool {
	return this.CallCount() > 0
}

func (this *MockController) Restore() {
	this.targetFunc.Set(this.originalFunc)
}
