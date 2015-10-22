package funcmock

// import . "github.com/onsi/gomega"

var empty = func(s string) int {
	return 0
}

var emptyMock = Mock(&empty)

func init() {
	empty("")
	empty("")
	empty("empty")
}

/*
func TestCallReverseOrder(t *testing.T) {
	RegisterTestingT(t)
	Expect(emptyMock.NthCall(0)).To(Equal(emptyMock.NthCall(-3)))
	Expect(emptyMock.NthCall(1)).To(Equal(emptyMock.NthCall(-2)))
	Expect(emptyMock.NthCall(2)).To(Equal(emptyMock.NthCall(-1)))

	Expect(emptyMock.NthCall(0)).NotTo(Equal(emptyMock.NthCall(1)))
	Expect(emptyMock.NthCall(0)).NotTo(Equal(emptyMock.NthCall(-2)))
}
*/
