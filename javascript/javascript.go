package javascript

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"os"
	"time"
)

func Asd() {
	vm := otto.New()

	// 定义一个 JavaScript 函数
	//va, err := vm.Run(`
	//   function add(a, b) {
	//       return a + b;
	//   }
	//`)
	va, err := vm.Run(`
	     console.log('sdsd');
	`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	integer, err := va.ToInteger()
	fmt.Println(integer + 2) // 输出: 8

	// 调用 JavaScript 函数
	value, err := vm.Call("add", nil, 5, 2)
	if err != nil {
		fmt.Println("Error calling function:", err)
		return
	}
	fmt.Println(value)
	toString, err := value.ToString()
	fmt.Println(toString)
	result, _ := value.ToInteger()
	json, err := value.MarshalJSON()
	fmt.Println(string(json))
	fmt.Println(result) // 输出: 8
}

var (
	halt      = errors.New("Stahp")
	typeError = errors.New("Unsupported type")
)

type Option func(c *SandboxJs)

func WithTimeout(timeout time.Duration) Option {
	return func(c *SandboxJs) {
		c.timeout = timeout
	}
}

func WithFuncName(funcName string) Option {
	return func(c *SandboxJs) {
		c.funcName = funcName
	}
}

type SandboxJs struct {
	otto    *otto.Otto
	timeout time.Duration

	funcName  string
	resultNum int
	args      []interface{}
	result    interface{}
}

func NewSandboxJs(ops ...Option) *SandboxJs {
	s := new(SandboxJs)
	s.otto = otto.New()
	for _, op := range ops {
		op(s)
	}
	return s
}

func (s *SandboxJs) ParseArgs(args ...interface{}) error {
	s.args = args
	return nil
}

func (s *SandboxJs) RunScript(code string) error {
	//是否开启超时
	s.timeoutInterrupt()
	if _, err := s.otto.Run(code); err != nil {
		return err
	}
	if len(s.args) < 1 {
		return nil
	}
	if len(s.funcName) == 0 {
		return nil
	}
	value, err := s.otto.Call(s.funcName, nil, s.args...)
	if err != nil {
		return err
	}
	resp, err := toGoType(value)
	if err != nil {
		return err
	}

	s.result = resp
	return nil
}

func (s *SandboxJs) GetResult() interface{} {
	return s.result
}

func (s *SandboxJs) timeoutInterrupt() {
	if s.timeout == 0 {
		return
	}
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == halt {
				fmt.Fprintf(os.Stderr, "Some code took to long! Stopping after: %v\n", duration)
				return
			}
			panic(caught) // Something else happened, repanic!
		}
		fmt.Fprintf(os.Stderr, "Ran code successfully: %v\n", duration)
	}()

	s.otto.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	watchdogCleanup := make(chan struct{})
	defer close(watchdogCleanup)

	go func() {
		select {
		case <-time.After(s.timeout): // Stop after two seconds
			s.otto.Interrupt <- func() {
				panic(halt)
			}
		case <-watchdogCleanup:
		}
		close(s.otto.Interrupt)
	}()
}

func toGoType(value otto.Value) (interface{}, error) {
	if value.IsBoolean() {
		return value.ToBoolean()
	}
	if value.IsString() {
		return value.ToString()
	}
	if value.IsNaN() {
		return value.ToFloat()
	}
	if value.IsNumber() {
		return value.ToInteger()
	}

	if value.IsNull() {
		return nil, nil
	}

	if value.IsObject() {
		return value.MarshalJSON()
	}
	return nil, typeError
}
