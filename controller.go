package funcmock

import (
	"fmt"
	"reflect"
	"sync"
)

type call struct {
	params      []interface{}
	returns     []interface{}
	specReturns []reflect.Value
}

func (this *call) recordParams(vals []reflect.Value) {
	this.params = make([]interface{}, len(vals))
	for idx, val := range vals {
		this.params[idx] = val.Interface()
	}
}

type MockController struct {
	target, original reflect.Value
	calls            []call
	callCount        int
	defaultReturns   []reflect.Value
	lock             sync.Mutex
	behaviour        *reflect.Value
	preRecord        *reflect.Value
	preReturn        *reflect.Value
}

type callHandle struct {
	controller *MockController
	calln      int
}

func (this *MockController) getCall(calln int) *call {
	if len(this.calls) <= calln {
		//this growth factor guarantees amortized O(1) time for insertions
		this.calls = append(this.calls, make([]call, calln+1)...)
	}
	return &this.calls[calln]
}

func (this *MockController) SetDefaultReturn(rets ...interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.defaultReturns != nil {
		panic("Can only call SetDefaultReturn once")
	}
	this.defaultReturns = this.validateReturns(rets)
}

func (this *MockController) Call(params []reflect.Value) (rets []reflect.Value) {
	this.lock.Lock()
	defer this.lock.Unlock()

	call := this.getCall(this.callCount)
	this.callCount++

	if this.preRecord != nil {
		params = this.preRecord.Call(params)
	}
	call.params = valueSliceToInterfaces(params)

	if this.behaviour != nil {
		rets = this.behaviour.Call(params)
	} else {
		if call.specReturns == nil {
			if this.defaultReturns == nil {
				targetFnType := this.target.Type()
				nouts := targetFnType.NumOut()
				for i := 0; i < nouts; i++ {
					this.defaultReturns = append(this.defaultReturns, reflect.Zero(targetFnType.Out(i)))
				}
			}
			rets = this.defaultReturns
		} else {
			rets = call.specReturns
		}
	}

	if this.preReturn != nil {
		rets = this.preReturn.Call(rets)
	}
	call.returns = valueSliceToInterfaces(rets)

	return
}

func (this *MockController) CallCount() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return int(this.callCount)
}
func (this *MockController) Called() bool {
	return this.CallCount() > 0
}

func (this *MockController) NthCall(calln int) callHandle {
	if calln < 0 {
		panic("NthCall called with negative index")
	}
	return callHandle{
		controller: this,
		calln:      calln,
	}
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
	this.controller.getCall(this.calln).specReturns = this.controller.validateReturns(rets)
}

func (this *MockController) NthParams(paramn int) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := reflect.MakeSlice(reflect.SliceOf(this.target.Type().In(paramn)), int(this.callCount), int(this.callCount))
	for idx, call := range this.calls[:this.callCount] {
		val.Index(idx).Set(reflect.ValueOf(call.params[paramn]))
	}
	return val.Interface()
}

func (this *MockController) NthReturns(paramn int) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := reflect.MakeSlice(reflect.SliceOf(this.target.Type().Out(paramn)), int(this.callCount), int(this.callCount))
	for idx, call := range this.calls[:this.callCount] {
		val.Index(idx).Set(reflect.ValueOf(call.returns[paramn]))
	}
	return val.Interface()
}

func (this *MockController) SetBehaviour(fn interface{}) *MockController {
	this.lock.Lock()
	defer this.lock.Unlock()

	fnval := reflect.ValueOf(fn)
	fntype := fnval.Type()
	if tgtype := this.target.Type(); fntype != tgtype {
		panic(fmt.Sprintf("MockController.WithBehaviour: behaviour function has invalid type '%s', requires '%s'", fntype.String(), tgtype.String()))
	}
	this.behaviour = &fnval

	return this
}

func (this *MockController) SetPreRecord(fn interface{}) *MockController {
	this.lock.Lock()
	defer this.lock.Unlock()

	fnval := reflect.ValueOf(fn)
	fntype := fnval.Type()

	tgtype := this.target.Type()
	intypes := make([]reflect.Type, 0, tgtype.NumIn())
	for i := 0; i < tgtype.NumIn(); i++ {
		intypes = append(intypes, tgtype.In(i))
	}
	reqtype := reflect.FuncOf(intypes, intypes, tgtype.IsVariadic())
	if fntype != reqtype {
		panic(fmt.Sprintf("MockController.SetPreRecord: provided function has invalid type '%s', requires '%s'", fntype.String(), reqtype.String()))
	}
	this.preRecord = &fnval

	return this
}

func (this *MockController) SetPreReturn(fn interface{}) *MockController {
	this.lock.Lock()
	defer this.lock.Unlock()

	fnval := reflect.ValueOf(fn)
	fntype := fnval.Type()

	tgtype := this.target.Type()
	outtypes := make([]reflect.Type, 0, tgtype.NumOut())
	for i := 0; i < tgtype.NumOut(); i++ {
		outtypes = append(outtypes, tgtype.Out(i))
	}
	reqtype := reflect.FuncOf(outtypes, outtypes, false)
	if fntype != reqtype {
		panic(fmt.Sprintf("MockController.SetPreReturn: provided function has invalid type '%s', requires '%s'", fntype.String(), reqtype.String()))
	}
	this.preReturn = &fnval

	return this
}

func (this *MockController) Restore() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.target.Set(this.original)
}

func (this *MockController) validateReturns(ins []interface{}) []reflect.Value {
	var idx int
	defer func() {
		if p := recover(); p != nil {
			panic(fmt.Sprintf("MockController.validateReturns: %d:th return value, %s", idx, p.(string)))
		}
	}()
	rets := make([]reflect.Value, len(ins))
	fntype := this.target.Type()
	insLen := len(ins)
	for idx = 0; idx < insLen; idx++ {
		val := ins[idx]
		typ := fntype.Out(idx)
		if val == nil {
			switch typ.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan, reflect.Interface:
				rets[idx] = reflect.Zero(typ)
			default:
				panic(fmt.Sprintf("type '%s' is not nillable", typ))
			}
		} else {
			rets[idx] = reflect.ValueOf(val).Convert(typ)
		}
	}
	return rets
}

func valueSliceToInterfaces(vals []reflect.Value) []interface{} {
	ret := make([]interface{}, len(vals))
	for idx, val := range vals {
		ret[idx] = val.Interface()
	}
	return ret
}
