package funcmock

import (
	"reflect"
	"sync"
)

func Mock(targetFnPtr interface{}) *MockController {
	if argType := reflect.TypeOf(targetFnPtr); argType.Kind() != reflect.Ptr || argType.Elem().Kind() != reflect.Func {
		panic("invalid argument to Mock! must be pointer")
	}
	targetFn := reflect.ValueOf(targetFnPtr).Elem()
	targetFnType := targetFn.Type()

	controller := &MockController{
		calls:     []call{},
		callCount: 0,
		lock:      sync.Mutex{},
		original:  reflect.ValueOf(targetFn.Interface()),
		target:    targetFn,
	}

	mockFn := reflect.MakeFunc(
		targetFnType,
		func(inValues []reflect.Value) []reflect.Value {
			return controller.Call(inValues)
		},
	)

	targetFn.Set(mockFn)

	return controller
}
