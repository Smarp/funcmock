package funcmock

import "reflect"

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
