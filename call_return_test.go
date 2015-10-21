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
