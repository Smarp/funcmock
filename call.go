package funcmock

type call struct {

}

func (this *call) ParamNth(nth int) interface{} {
	return nil
}

func (this *call) Return(args ...interface{}) *call {
	return this
}
