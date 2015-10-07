package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

var idempotent = func(i int) int {
	return i
}

var idempotentMock = Mock(&idempotent)

func TestCallCounter0(*testing.T) {
	Expect(idempotentMock.CallCounter()).To(Equal(0))
}

func TestCalledFalse(*testing.T) {
	Expect(idempotentMock.Called()).To(BeFalse())
}
