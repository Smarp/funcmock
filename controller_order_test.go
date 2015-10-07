package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

var empty = func(s string) int {
	return 0
}

var emptyMock = Mock(&empty)

func init() {
	empty("")
	empty("")
	empty("empty")
}

func TestCallReverseOrder(*testing.T) {
	Expect(emptyMock.CallNth(0)).To(Equal(emptyMock.CallNth(-3)))
	Expect(emptyMock.CallNth(1)).To(Equal(emptyMock.CallNth(-2)))
	Expect(emptyMock.CallNth(2)).To(Equal(emptyMock.CallNth(-1)))

	Expect(emptyMock.CallNth(0)).NotTo(Equal(emptyMock.CallNth(1)))
	Expect(emptyMock.CallNth(0)).NotTo(Equal(emptyMock.CallNth(-2)))
}
