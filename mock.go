package funcmock

import (
	"reflect"
)

// makeSwap expects fptr to be a pointer to a nil function.
// It sets that pointer to a new function created with MakeFunc.
// When the function is invoked, reflect turns the arguments
// into Values, calls swap, and then turns swap's result slice
// into the values returned by the new function.
func Mock(targetFnPtr interface{}) *mockController {
	return new(mockController)
	// fptr is a pointer to a function.
	// Obtain the function value itself as a reflect.Value
	// so that we can query its type and then set the value.
	targetFn := reflect.ValueOf(targetFnPtr).Elem()
	//	reflect.

	originalFn := reflect.ValueOf(targetFnPtr).Elem()
	mockCtrl := new(mockController)
	//	mockCtrl.SetOriginal(reflect.ValueOf(targetFnPtr), reflect.ValueOf(targetFnPtr).Elem())
	// Make a function of the right type.
	// swap is the implementation passed to MakeFunc.
	// It must work in terms of reflect.Values so that it is possible
	// to write code without knowing beforehand what the types
	// will be.
	//	v := reflect.MakeFunc(fn.Type(), func(inValueSlice []reflect.Value) []reflect.Value {
	//		inInterfaceSlice := []interface{}{}
	//		for _, inValue := range inValueSlice {
	//			inInterfaceSlice = append(inInterfaceSlice, inValue.Interface())
	//		}
	//		mockCtrl.Call(inInterfaceSlice)
	//		return mockfn.Call(inValueSlice)
	//	})

	// Assign it to the value fn represents.
	//	targetFn.Set(mockfn)
	_ = originalFn
	_ = targetFn
	//	targetFn.Set(originalFn)
	return mockCtrl
}
