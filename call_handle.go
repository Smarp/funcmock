package funcmock

import "fmt"

type callHandle struct {
	controller *MockController
	calln      int
}

func (this callHandle) NthParam(paramn int) interface{} {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	if this.controller.callCount <= this.calln {
		panic(fmt.Sprintf("%dth call has not been made", this.calln))
	}
	call := &this.controller.calls[this.calln]
	return call.params[paramn]
}

func (this callHandle) NthReturn(ret int) interface{} {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	if this.controller.callCount <= this.calln {
		panic(fmt.Sprintf("%dth call has not been made", this.calln))
	}
	return this.controller.calls[this.calln].returns[ret]
}

func (this callHandle) SetReturn(rets ...interface{}) {
	this.controller.lock.Lock()
	defer this.controller.lock.Unlock()
	this.controller.getCall(this.calln).specReturns = this.controller.sanitizeReturns(rets)
}
