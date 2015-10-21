package funcmock

import (
	"testing"

	. "github.com/onsi/gomega"
)

var idempotent = func(i int) int {
	return i
}

var idempotentMock = Mock(&idempotent)

func TestCallCounter0(t *testing.T) {
	RegisterTestingT(t)
	Expect(idempotentMock.CallCounter()).To(Equal(0))
}

func TestCalledFalse(t *testing.T) {
	RegisterTestingT(t)
	Expect(idempotentMock.Called()).To(BeFalse())
}
