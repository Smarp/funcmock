package funcmock

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var mockTestFuncArgs []string
var mockTestFunc = func(msg string) {
	mockTestFuncArgs = append(mockTestFuncArgs, msg)
}

var recordingValuesTestFunc = func(a int, b string) (int, string) {
	Fail("This function should not be called")
	return 0, ""
}

var interfaceReturnTestFunc = func() error {
	Fail("This function should not be called")
	return nil
}

var stressTestFunc = func(a int) {
	Fail("This function should not be called")
}

var setBehaviourTestFunc = func(a, b int) int {
	Fail("This function should not be called")
	return 0
}

type preRecordTestStruct struct {
	field int
}

var preRecordTestFunc = func(ts []preRecordTestStruct) *preRecordTestStruct {
	Fail("This function should not be called")
	return nil
}

var _ = Describe("Better Mock Test", func() {
	Context("mocking/restoring", func() {

		var mock *MockController

		BeforeEach(func() {
			mockTestFuncArgs = make([]string, 0)
		})
		It("should mock and restore the function correctly", func() {
			mockTestFunc("call1")
			mock = Mock(&mockTestFunc)
			mockTestFunc("call2")
			mock.Restore()
			mockTestFunc("call3")

			Expect(mockTestFuncArgs).To(Equal([]string{"call1", "call3"}))
		})
		It("should panic if Mock is called with non-pointer argument", func() {
			var err interface{}
			func() {
				defer func() {
					err = recover()
				}()
				mock = Mock(mockTestFunc)
			}()
			Expect(err).To(Equal("invalid argument to Mock! must be pointer to function"))
		})
		It("should panic if Mock is called with non-function type", func() {
			var err interface{}
			func() {
				defer func() {
					err = recover()
				}()
				mock = Mock(&struct{ field int }{})
			}()
			Expect(err).To(Equal("invalid argument to Mock! must be pointer to function"))
		})
	})
	Context("recording values", func() {

		var mock *MockController

		BeforeEach(func() {
			mock = Mock(&recordingValuesTestFunc)
		})
		AfterEach(func() {
			mock.Restore()
		})
		It("should correctly record the number of calls", func() {
			for i := 0; i < 5; i++ {
				recordingValuesTestFunc(0, "")
			}
			Expect(mock.CallCount()).To(Equal(5))
			for i := 0; i < 5; i++ {
				Expect(mock.NthCall(0).Called()).To(BeTrue())
			}
			Expect(mock.NthCall(5).Called()).To(BeFalse())
		})
		It("should record and retrieve the call parameters", func() {
			for i := 0; i < 25; i++ {
				recordingValuesTestFunc(i%5, fmt.Sprintf("%d", i/5))
			}
			for i := 0; i < 25; i++ {
				Expect(mock.NthCall(i).NthParam(0)).To(Equal(i % 5))
				Expect(mock.NthCall(i).NthParam(1)).To(Equal(fmt.Sprintf("%d", i/5)))
			}
		})
		It("should be able to retrieve all call parameters with NthParams", func() {
			for i := 0; i < 5; i++ {
				recordingValuesTestFunc(i*2, fmt.Sprintf("%d", i))
			}
			params1, ok := mock.NthParams(0).([]int)
			Expect(ok).To(BeTrue())
			params2, ok := mock.NthParams(1).([]string)
			Expect(ok).To(BeTrue())
			Expect(params1).To(Equal([]int{0, 2, 4, 6, 8}))
			Expect(params2).To(Equal([]string{"0", "1", "2", "3", "4"}))
		})
		It("should have zero default return values", func() {
			a, b := recordingValuesTestFunc(0, "")
			Expect(a).To(Equal(0))
			Expect(b).To(Equal(""))
		})
		It("should support setting default return values", func() {
			mock.SetDefaultReturn(2, "test")
			a, b := recordingValuesTestFunc(0, "")
			Expect(a).To(Equal(2))
			Expect(b).To(Equal("test"))
		})
		It("should panic if SetDefaultReturn is called with invalid type", func() {
			var err error
			func() {
				defer func() {
					err = errors.New(recover().(string))
				}()
				mock.SetDefaultReturn("test", 2)
			}()
			Expect(err).To(Equal(errors.New("MockController.sanitizeReturns: 0:th return value, reflect.Value.Convert: value of type string cannot be converted to type int")))
		})
		Context("per call return values", func() {
			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					mock.NthCall(i).SetReturn(i, fmt.Sprintf("%d", i*2))
				}
			})
			It("should support setting return values for specific calls", func() {
				for i := 0; i < 5; i++ {
					a, b := recordingValuesTestFunc(0, "")
					Expect(a).To(Equal(i))
					Expect(b).To(Equal(fmt.Sprintf("%d", i*2)))
				}
				a, b := recordingValuesTestFunc(0, "")
				Expect(a).To(Equal(0))
				Expect(b).To(Equal(""))
			})
			It("should panic if SetReturn is called with invalid types", func() {
				var err error
				func() {
					defer func() {
						err = errors.New(recover().(string))
					}()
					mock.NthCall(0).SetReturn("test", 2)
				}()
				Expect(err).To(Equal(errors.New("MockController.sanitizeReturns: 0:th return value, reflect.Value.Convert: value of type string cannot be converted to type int")))
			})
			It("should record the returned values", func() {
				for i := 0; i < 6; i++ {
					recordingValuesTestFunc(0, "")
				}
				for i := 0; i < 5; i++ {
					Expect(mock.NthCall(i).NthReturn(0)).To(Equal(i))
					Expect(mock.NthCall(i).NthReturn(1)).To(Equal(fmt.Sprintf("%d", i*2)))
				}
				Expect(mock.NthCall(5).NthReturn(0)).To(Equal(0))
				Expect(mock.NthCall(5).NthReturn(1)).To(Equal(""))
			})
			It("should be able to retrieve all returned values with NthReturns", func() {
				for i := 0; i < 6; i++ {
					recordingValuesTestFunc(0, "")
				}
				returns1, ok := mock.NthReturns(0).([]int)
				Expect(ok).To(BeTrue())
				returns2, ok := mock.NthReturns(1).([]string)
				Expect(ok).To(BeTrue())
				Expect(returns1).To(Equal([]int{0, 1, 2, 3, 4, 0}))
				Expect(returns2).To(Equal([]string{"0", "2", "4", "6", "8", ""}))
			})
		})
	})
	Context("interface returns", func() {

		var mock *MockController

		BeforeEach(func() {
			mock = Mock(&interfaceReturnTestFunc)
		})
		AfterEach(func() {
			mock.Restore()
		})
		It("should return nil as a default value", func() {
			ret := interfaceReturnTestFunc()
			Expect(ret).To(BeNil())
		})
		It("should allow setting nil default return for functions returning interfaces", func() {
			mock.SetDefaultReturn(nil)
			ret := interfaceReturnTestFunc()
			Expect(ret).To(BeNil())
		})
		It("should allow setting return values which satisfy the interface", func() {
			mock.SetDefaultReturn(errors.New("test error"))
			ret := interfaceReturnTestFunc()
			Expect(ret).To(Equal(errors.New("test error")))
		})
	})
	Context("stress test", func() {

		var mock *MockController

		BeforeEach(func() {
			mock = Mock(&stressTestFunc)
		})
		AfterEach(func() {
			mock.Restore()
		})
		It("should support calls from multiple goroutines without race conditions or deadlocks", func() {
			goroutines := 20
			callsPerGoroutine := 1000
			totalCalls := goroutines * callsPerGoroutine
			callParams := make([]int, totalCalls)
			for i := 0; i < totalCalls; i++ {
				callParams[i] = i
			}
			for i := 0; i < goroutines; i++ {
				go func(slice []int) {
					for j := 0; j < callsPerGoroutine; j++ {
						stressTestFunc(slice[j])
					}
				}(callParams[i*callsPerGoroutine : (i+1)*callsPerGoroutine])
			}
			Eventually(func() int {
				return mock.CallCount()
			}, 1.0, 0.1).Should(Equal(totalCalls))
			// this is done in this way instead of using gomega.ConsistOf since it's too slow
			recordedParams := mock.NthParams(0).([]int)
			for _, val := range recordedParams {
				if val < 0 || val >= totalCalls || callParams[val] != val {
					Fail("Recorded params are invalid")
				} else {
					callParams[val] = -1
				}
			}
		})
	})
	Context("overriding behaviour", func() {

		var mock *MockController

		Context("SetBehaviour", func() {
			BeforeEach(func() {
				mock = Mock(&setBehaviourTestFunc)
			})
			AfterEach(func() {
				mock.Restore()
			})
			It("should fail if SetBehaviour is called with invalid function type", func() {
				var err error
				func() {
					defer func() {
						err = errors.New(recover().(string))
					}()
					mock.SetBehaviour(func(x int) int { return x })
				}()
				Expect(err).To(Equal(errors.New("MockController.WithBehaviour: behaviour function has invalid type 'func(int) int', requires 'func(int, int) int'")))
			})
			It("should allow overriding mock function behaviour", func() {
				type pair struct{ a, b int }
				args := make([]pair, 0)
				mock.SetBehaviour(func(a, b int) int {
					args = append(args, pair{a, b})
					return a * b
				})
				Expect(setBehaviourTestFunc(1, 2)).To(Equal(2))
				Expect(setBehaviourTestFunc(3, 4)).To(Equal(12))
				Expect(setBehaviourTestFunc(5, 6)).To(Equal(30))
				Expect(args).To(Equal([]pair{{1, 2}, {3, 4}, {5, 6}}))
			})
			It("should store call info with overridden behaviour", func() {
				mock.SetBehaviour(func(a, b int) int {
					return a * b
				})
				Expect(setBehaviourTestFunc(1, 2)).To(Equal(2))
				Expect(setBehaviourTestFunc(3, 4)).To(Equal(12))
				Expect(setBehaviourTestFunc(5, 6)).To(Equal(30))

				Expect(mock.CallCount()).To(Equal(3))

				Expect(mock.NthCall(0).NthParam(0)).To(Equal(1))
				Expect(mock.NthCall(0).NthParam(1)).To(Equal(2))
				Expect(mock.NthCall(1).NthParam(0)).To(Equal(3))
				Expect(mock.NthCall(1).NthParam(1)).To(Equal(4))
				Expect(mock.NthCall(2).NthParam(0)).To(Equal(5))
				Expect(mock.NthCall(2).NthParam(1)).To(Equal(6))

				Expect(mock.NthCall(0).NthReturn(0)).To(Equal(2))
				Expect(mock.NthCall(1).NthReturn(0)).To(Equal(12))
				Expect(mock.NthCall(2).NthReturn(0)).To(Equal(30))
			})
		})
		Context("SetPreRecord and SetPreReturn", func() {
			BeforeEach(func() {
				mock = Mock(&preRecordTestFunc)
			})
			AfterEach(func() {
				mock.Restore()
			})
			It("should fail if SetPreRecord is called with invalid function type", func() {
				var err error
				func() {
					defer func() {
						err = errors.New(recover().(string))
					}()
					mock.SetPreRecord(func(ts []preRecordTestStruct) *preRecordTestStruct { return nil })
				}()
				Expect(err).To(Equal(errors.New(
					"MockController.SetPreRecord: " +
						"provided function has invalid type 'func([]funcmock.preRecordTestStruct) *funcmock.preRecordTestStruct', " +
						"requires 'func([]funcmock.preRecordTestStruct) []funcmock.preRecordTestStruct'",
				)))
			})
			It("should allow changing the recording behaviour", func() {
				mock.SetPreRecord(func(ts []preRecordTestStruct) []preRecordTestStruct {
					ret := make([]preRecordTestStruct, len(ts))
					copy(ret, ts)
					return ret
				})
				rets := []preRecordTestStruct{{field: 1}, {field: 3}}
				preRecordTestFunc(rets)
				rets[0].field = 2
				rets[1].field = 4
				preRecordTestFunc(rets)
				Expect(mock.NthCall(1).NthParam(0)).To(Equal([]preRecordTestStruct{{field: 2}, {field: 4}}))
				Expect(mock.NthCall(0).NthParam(0)).To(Equal([]preRecordTestStruct{{field: 1}, {field: 3}}))
			})
			It("should fail if SetPreReturn is called with invalid function type", func() {
				var err error
				func() {
					defer func() {
						err = errors.New(recover().(string))
					}()
					mock.SetPreReturn(func(ts []preRecordTestStruct) *preRecordTestStruct { return nil })
				}()
				Expect(err).To(Equal(errors.New(
					"MockController.SetPreReturn: " +
						"provided function has invalid type 'func([]funcmock.preRecordTestStruct) *funcmock.preRecordTestStruct', " +
						"requires 'func(*funcmock.preRecordTestStruct) *funcmock.preRecordTestStruct'",
				)))
			})
			It("should return the values returned by the function passed to SetPreReturn from the mocked function", func() {
				mock.SetPreReturn(func(ts *preRecordTestStruct) *preRecordTestStruct {
					temp := *ts
					return &temp
				})
				mock.SetDefaultReturn(&preRecordTestStruct{field: 1})
				var rets [3]*preRecordTestStruct
				for i := 0; i < 3; i++ {
					st := preRecordTestFunc(nil)
					st.field++
					rets[i] = st
				}
				Expect(rets[:]).To(Equal([]*preRecordTestStruct{
					&preRecordTestStruct{field: 2},
					&preRecordTestStruct{field: 2},
					&preRecordTestStruct{field: 2},
				}))
			})
			It("should record the values returned by the function passed to SetPreReturn", func() {
				mock.SetPreReturn(func(ts *preRecordTestStruct) *preRecordTestStruct {
					temp := *ts
					return &temp
				})
				mock.SetDefaultReturn(&preRecordTestStruct{field: 1})
				for i := 0; i < 3; i++ {
					st := preRecordTestFunc(nil)
					st.field++
				}
				for i := 0; i < 3; i++ {
					Expect(mock.NthCall(i).NthReturn(0)).To(PointTo(Equal(preRecordTestStruct{field: 2})))
				}
			})
		})
	})
})
