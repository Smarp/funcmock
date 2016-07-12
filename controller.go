package funcmock

import (
	"reflect"
	"sync"
)

type MockController struct {
	originalFunc reflect.Value
	targetFunc   reflect.Value
	// A lock for callStack
	callsMutex *sync.Mutex
	// we need map, not slice, to set call before it is called
	callStack map[int]*call
	// A lock for counter
	counterMutex *sync.Mutex
	counter      int
	// the default call which shall be used for mint calls
	defaultYield []reflect.Value

	// Flag indicating the default return has been set
	yieldSet bool
}

func (this *MockController) CallCount() (count int) {
	this.counterMutex.Lock()
	count = this.counter
	this.counterMutex.Unlock()
	return
}

func (this *MockController) NthCall(nth int) (theCall *call) {
	this.callsMutex.Lock()
	theCall, ok := this.callStack[nth]
	if !ok {
		theCall = &call{
			paramMutex: &sync.Mutex{},
			called:     false,
		}
		this.callStack[nth] = theCall
	}
	this.callsMutex.Unlock()
	return
}

func (this *MockController) incrementCounter() (count int) {
	this.counterMutex.Lock()
	this.counter++
	count = this.counter
	this.counterMutex.Unlock()
	return
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
	this.callsMutex.Lock()
	this.callStack[index] = theCall
	this.callsMutex.Unlock()
	return
}

func (this *MockController) Called() bool {
	return this.CallCount() > 0
}

func (this *MockController) Restore() {
	this.targetFunc.Set(this.originalFunc)
}
