package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCall0thReturn(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}

	var swapMock = Mock(&swap)
	swapMock.CallNth(0).SetReturn(8, 9)
	Expect(swapMock.CallNth(0)).NotTo(BeNil())
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
	swapMock.CallNth(2).SetReturn(8, 9)
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
	swapMock.SetDefaultReturn("three", "four")
	var panicMessage interface{}
	func() {
		defer func() {
			panicMessage = recover()
		}()
		_, _ = swap(2, 2)
	}()
	Expect(panicMessage).To(Equal("reflect: function created by MakeFunc using closure returned wrong type: have string for int"))
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
