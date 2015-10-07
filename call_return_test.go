package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCall0thReturn(*testing.T) {
	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.CallNth(0).Return(8, 9)
	v8, v9 := swap(2, 2)
	Expect(v8).To(Equal(8))
	Expect(v9).To(Equal(9))
}

func TestCallLastReturn(*testing.T) {
	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.CallNth(2).Return(8, 9)
	_ = swap(2, 2)
	_ = swap(2, 2)
	v8, v9 := swap(2, 2)
	Expect(v8).To(Equal(8))
	Expect(v9).To(Equal(9))
}

func TestCallOtherReturn(*testing.T) {
	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.CallNth(0).Return(2, 3)
	swapMock.CallOther().Return(12, 13)
	v2, v3 := swap(5, 6)
	Expect(v2).To(Equal(2))
	Expect(v3).To(Equal(3))
	v12, v13 := swap(5, 6)
	Expect(v12).To(Equal(12))
	Expect(v13).To(Equal(13))
}
