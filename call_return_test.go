package funcmock

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCall0thReturn(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}

	var swapMock = Mock(&swap)
	swapMock.NthCall(0).SetReturn(8, 9)
	Expect(swapMock.NthCall(0)).NotTo(BeNil())
	v8, v9 := swap(2, 2)
	Expect(v8).To(Equal(8))
	Expect(v9).To(Equal(9))
}

func TestCallLastReturn(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.NthCall(2).SetReturn(8, 9)
	_, _ = swap(2, 2)
	_, _ = swap(2, 2)
	v8, v9 := swap(2, 2)
	Expect(v8).To(Equal(8))
	Expect(v9).To(Equal(9))
}

func TestSetDefaultReturnWrongType(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}
	_, _ = swap(1, 1)
	var swapMock = Mock(&swap)
	var panicMessage interface{}
	func() {
		defer func() {
			panicMessage = recover()
		}()
		swapMock.SetDefaultReturn("three", "four")
	}()
	Expect(panicMessage).To(Equal("reflect.Value.Convert: value of type string cannot be converted to type int"))
}

func TestSetDefaultReturnNil(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (*int, *int) {
		return &j, &i
	}
	_, _ = swap(1, 1)
	var swapMock = Mock(&swap)
	swapMock.SetDefaultReturn(nil, nil)
	Expect(func() { _, _ = swap(2, 2) }).NotTo(Panic())
	v8, v9 := swap(2, 2)
	Expect(v8).To(BeNil())
	Expect(v9).To(BeNil())
}

func TestSetReturnNil(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (*int, *int) {
		return &j, &i
	}
	_, _ = swap(1, 1)
	var swapMock = Mock(&swap)
	swapMock.NthCall(0).SetReturn(nil, nil)
	Expect(func() { _, _ = swap(2, 2) }).NotTo(Panic())
	v8, v9 := swap(2, 2)
	Expect(v8).To(BeNil())
	Expect(v9).To(BeNil())
}
func TestSetReturnWithErrorReturnParam(t *testing.T) {
	RegisterTestingT(t)

	var funcToTest = func(i int, j int) (err error) {
		return err
	}
	var swapMock = Mock(&funcToTest)
	swapMock.NthCall(0).SetReturn(errors.New("message"))
	Expect(func() { _ = funcToTest(2, 2) }).NotTo(Panic())
	swapMock.NthCall(1).SetReturn(errors.New("message"))
	err := funcToTest(2, 2)
	Expect(err.Error()).To(Equal("message"))
}
