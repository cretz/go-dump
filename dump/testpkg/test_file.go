package testpkg

import (
	"io"
	"strconv"
)

var (
	varIntArr        = [4]int{1, 2, 3, 4}
	varInt           = 5
	varIntChan       chan int
	varIntChanSend   chan<- int
	varIntChanRecv   <-chan int
	varAnonInterface interface {
		io.Reader
		Foo() string
	}
	varMap           map[string]string
	varStringPointer *string
	varStringSlice   []string
	varAnonStruct    struct {
		io.Reader
		Foo string
	}
	varFunc func(string) error
)

const (
	constBool                      = true
	constString                    = "str"
	constInt            int        = 5
	constIntUnknown                = 5
	constFloat          float32    = 1.2
	constFloatUnknown              = 1.2
	constComplex        complex128 = (0.0 + 1.0i)
	constComplexUnknown            = (0.0 + 1.0i)
)

func TopLevelFunc(foo int) string {
	return strconv.Itoa(foo)
}
