package funcmock

import "reflect"

func Mock(targetFnPtr interface{}) (controller *MockController) {

	targetFn := reflect.ValueOf(targetFnPtr).Elem()
	controller = &MockController{
		counter:   make(chan int),
		callStack: make(chan map[int]*call),
	}
	go func() {
		controller.counter <- 0
		controller.callStack <- make(map[int]*call)
	}()
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
			controller.incrementCounter()
			theCall, ok := controller.callStack[controller.CallCount()-1]
			if ok == false {
				theCall = &call{
					param: make(chan []interface{}),
				}
				controller.add(theCall)
			}

			{
				// could be moved out
				var param []interface{}
				go func() { theCall.param <- param }()

			}
			theCall.called = true
			for i := 0; i < targetFnType.NumIn(); i++ {
				theCall.appendParam(inValueSlice[i].Interface())

			}
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
