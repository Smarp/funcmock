package funcmock

import (
	"sync"
	"testing"

	. "github.com/onsi/gomega"
)

func TestZeroValueInt(*testing.T) {
	reverse := func(i int) int {
		return -i
	}
	Expect(reverse(1)).To(Equal(-1))
	mockCtrl := Mock(&reverse)
	Expect(reverse(1)).To(Equal(0), "zero value")
	mockCtrl.Restore()
	Expect(reverse(1)).To(Equal(-1))
}

func TestZeroValueString(*testing.T) {
	prepend := func(i string) string {
		return "prefix" + i
	}
	Expect(prepend("body")).To(Equal("prefixbody"))
	mockCtrl := Mock(&prepend)
	Expect(prepend("body")).To(Equal(""), "zero value")
	mockCtrl.Restore()
	Expect(prepend("body")).To(Equal("prefixbody"))
}

func TestRaceCondition(*testing.T) {
	reverse := func(i int) int {
		return -i
	}
	mockCtrl := Mock(&reverse)
	var wg sync.WaitGroup
	const z = 1000
	wg.Add(z)
	for i := 0; i < z; i++ {
		go func() {
			defer wg.Done()
			reverse(i)
		}()
	}
	wg.Wait()
	Expect(mockCtrl.CallCounter()).To(Equal(z))
}
