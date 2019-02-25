package funcmock

import "fmt"

type callHandle struct {
	controller *MockController
	callIndex  int
}

func (this callHandle) Called() bool {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	return this.controller.callCount > this.callIndex
}

func (this callHandle) NthParam(paramn int) interface{} {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	if this.controller.callCount <= this.callIndex {
		panic(fmt.Sprintf("%dth call has not been made", this.callIndex))
	}
	call := &this.controller.calls[this.callIndex]
	return call.params[paramn]
}

func (this callHandle) NthReturn(ret int) interface{} {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	if this.controller.callCount <= this.callIndex {
		panic(fmt.Sprintf("%dth call has not been made", this.callIndex))
	}
	return this.controller.calls[this.callIndex].returns[ret]
}

func (this callHandle) SetReturn(rets ...interface{}) {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	this.controller.getCall(this.callIndex).specReturns = this.controller.sanitizeReturns(rets)
}
