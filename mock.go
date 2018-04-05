package funcmock

import (
	"reflect"
)

func Mock(targetFnPtr interface{}) (controller *MockController) {

	targetFn := reflect.ValueOf(targetFnPtr).Elem()
	controller = &MockController{
		counter:   make(chan int),
		callStack: make(chan map[int]*call),
	}
	go func() { controller.callStack <- make(map[int]*call) }()
	go func() {controller.counter <- 0 }()

	controller.targetFunc = targetFn
	targetFnType := targetFn.Type()
	numberOfOuts := targetFnType.NumOut()

	controller.originalFunc = reflect.ValueOf(targetFn.Interface())

	for i := 0; i < numberOfOuts; i++ {
		controller.defaultYield = append(controller.defaultYield,
			reflect.Zero(targetFnType.Out(i)))
	}

	mockFn := reflect.MakeFunc(targetFnType,
		func(inValueSlice []reflect.Value) (yield []reflect.Value) {
			callCount := controller.incrementCounter()
			theCall := controller.NthCall(callCount - 1)
			var param []interface{}
			for i := 0; i < targetFnType.NumIn(); i++ {
				param = append(param, inValueSlice[i].Interface())
			}
			theCall.updateParam(param)
			if numberOfOuts == len(theCall.yield) {
				// if user has set the return values the spit them out
				for i := 0; i < numberOfOuts; i++ {
					yield = append(yield,
						sanitizeReturn(
							targetFnType.Out(i),
							theCall.yield[i]))
				}
			} else {
				yield = controller.defaultYield
			}
			return yield
		},
	)

	targetFn.Set(mockFn)

	return controller
}
