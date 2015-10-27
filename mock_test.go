package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestZeroValueInt(t *testing.T) {
	RegisterTestingT(t)
	reverse := func(i int) int {
		return -i
	}
	Expect(reverse(1)).To(Equal(-1))
	mockCtrl := Mock(&reverse)
	Expect(reverse(1)).To(Equal(0), "zero value")
	mockCtrl.Restore()
	Expect(reverse(1)).To(Equal(-1))
}

func TestZeroValueString(t *testing.T) {
	RegisterTestingT(t)
	prepend := func(i string) string {
		return "prefix" + i
	}
	Expect(prepend("body")).To(Equal("prefixbody"))
	mockCtrl := Mock(&prepend)
	Expect(prepend("body")).To(Equal(""), "zero value")
	mockCtrl.Restore()
	Expect(prepend("body")).To(Equal("prefixbody"))
}

func TestSetDefaultReturn(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.NthCall(0).SetReturn(2, 3)
	swapMock.SetDefaultReturn(12, 13)
	v2, v3 := swap(5, 6)
	Expect(v2).To(Equal(2))
	Expect(v3).To(Equal(3))
	v12, v13 := swap(5, 6)
	Expect(v12).To(Equal(12))
	Expect(v13).To(Equal(13))
}

func TestCallSetDefaultReturnOnce(t *testing.T) {
	RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.NthCall(0).SetReturn(2, 3)
	Expect(func() { swapMock.SetDefaultReturn(1, 3) }).NotTo(Panic())
	Expect(func() { swapMock.SetDefaultReturn(1, 3) }).To(Panic())

}

/*
func TestMultiSetDefaultReturnCalls(t *testing.T){
RegisterTestingT(t)

	var swap = func(i, j int) (int, int) {
		return j, i
	}
	var swapMock = Mock(&swap)
	swapMock.NthCall(0).SetReturn(2, 3)
	swapMock.CallOther().SetReturn(12, 13)
	v1, v2 = swap(5, 6)
	v3, v4 = swap(14, 17)
	swapMock.CallOther().SetReturn(45, 56)
	v5, v6 = swap(67,89)
	Expect(v1).To(Equal(2))
	Expect(v2).To(Equal(3))
	Expect(v3).To(Equal(12))
	Expect(v4).To(Equal(13))
	Expect(v5).To(Equal(45))
	Expect(v6).To(Equal(56))

}
*/

func TestRaceCondition(t *testing.T) {
	RegisterTestingT(t)
	// t.SkipNow()
	testing.Benchmark(func(b *testing.B) {
		reverse := func(i int) int {
			return -i
		}
		// b.N = 100000
		mockCtrl := Mock(&reverse)
		// RunParallel will create GOMAXPROCS goroutines
		// and distribute work among them.
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				i := 1
				reverse(i)
			}
		})
		Expect(mockCtrl.CallCount()).To(Equal(b.N))
	})
}
