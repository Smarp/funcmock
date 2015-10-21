package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

var swap = func(i, j int) (int, int) {
	return j, i
}

var swapMock = Mock(&swap)

func init() {
	swap(3, -5)
	swap(-7, 11)
	swap(13, -17)
}

func TestCallCounter3(t *testing.T) {
	RegisterTestingT(t)
	Expect(swapMock.CallCounter()).To(Equal(3))
}

func TestCalledTrue(t *testing.T) {
	RegisterTestingT(t)
	Expect(swapMock.Called()).To(BeTrue())
}

func TestCall0thParams(t *testing.T) {
	RegisterTestingT(t)
	call0nth := swapMock.CallNth(0)
	Expect(call0nth).NotTo(BeNil())
	Expect(call0nth.ParamNth(0)).To(Equal(3))
	Expect(call0nth.ParamNth(1)).To(Equal(-5))
}

func TestCallLastParams(t *testing.T) {
	RegisterTestingT(t)
	callLast := swapMock.CallNth(2)
	Expect(callLast.ParamNth(0)).To(Equal(13))
	Expect(callLast.ParamNth(1)).To(Equal(-17))
}
