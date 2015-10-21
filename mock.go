package funcmock

import "reflect"

func Mock(targetFnPtr interface{}) (controller *mockController) {

	targetFn := reflect.ValueOf(targetFnPtr).Elem()
	controller = &mockController{
		callStack: make(map[int]*call),
		counter:   0,
	}
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
			theCall, ok := controller.callStack[controller.CallCounter()-1]
			if ok == false {
				theCall = new(call)
				controller.add(theCall)
			}
			theCall.called = true
			for i := 0; i < targetFnType.NumIn(); i++ {
				theCall.param = append(theCall.param, inValueSlice[i].Interface())

			}
			if numberOfOuts == len(theCall.yield) {
				// if user has set the return values the spit them out
				for i := 0; i < numberOfOuts; i++ {
					yield = append(yield, reflect.ValueOf(theCall.yield[i]))
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
