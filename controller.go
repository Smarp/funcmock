package funcmock

import (
	"reflect"
)

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
	var count int
	select {
	case count = <-this.counter:
		go func() { this.counter <- count }()

	}
	return count
}

func (this *MockController) NthCall(nth int) (theCall *call) {
	callStack := <-this.callStack
	theCall, ok := callStack[nth]
	if ok == false {
		theCall = &call{
			param:  make(chan []interface{}),
			called: false,
		}

	}

	go func() { this.callStack <- callStack }()
	this.addCallAt(theCall, nth)
	return theCall
}

func (this *MockController) incrementCounter() int {
	var count int
	select {
	case count = <-this.counter:
		count++
		go func() { this.counter <- count }()

	}
	return count
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

// func (this *MockController) getCallStack() map[int]*call {
// 	go func() { this.callStack <- callStack }()
// 	return callStack
// }

func (this *MockController) addCallAt(theCall *call, index int) {
	callStack := <-this.callStack
	callStack[index] = theCall
	go func() { this.callStack <- callStack }()
}

func (this *MockController) Called() bool {
	return this.CallCount() > 0
}

func (this *MockController) Restore() {
	this.targetFunc.Set(this.originalFunc)
}
