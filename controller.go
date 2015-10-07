package funcmock

import (
	"reflect"
)

type mockController struct {
	originalFunc   reflect.Value
	targetFunc     reflect.Value
	oFuncPtr       interface{}
	oFunc          interface{}
	callStackStack map[int]*call
}

func (this *mockController) CallCounter() int {
	return len(this.callStackStack)
}

func (this *mockController) CallNth(nth int) *call {
	return nil
}

func (this *mockController) CallOther() *call {
	return nil
}

func (this *mockController) Called() bool {
	return false
}

func (this *mockController) Restore() {
	this.targetFunc.Set(this.originalFunc)
}
