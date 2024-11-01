package lua

import (
	"context"
	"errors"
	"fmt"
	"github.com/mikez/go-vm/utils"
	lua "github.com/yuin/gopher-lua"
	"time"
)

func RunLua() {
	code :=
		`function calculate(a, b)
    		local sum = a + b
    		local difference = a - b
			return sum, difference
		end`
	L := lua.NewState()
	defer L.Close()
	//if err := L.DoString(`print("hello")`); err != nil {
	//	panic(err)
	//}
	if err := L.DoString(code); err != nil {
		panic(err)
	}
	// 调用 Lua 函数并获取返回值
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("calculate"),
		NRet:    2, // 期望两个返回值
		Protect: true,
	}, lua.LNumber(3), lua.LNumber(5)); err != nil {
		panic(err)
	}

	// 获取返回值
	ret1 := L.Get(-2) // 第一个返回值
	ret2 := L.Get(-1) // 第二个返回值
	L.Pop(2)          // 弹出返回值

	fmt.Printf("Result: %v\n", ret1)
	fmt.Printf("Result: %v\n", ret2)
}

var (
	TypeError  = errors.New("args type error")
	CloseError = errors.New("lua have been close error")
)

type Option func(c *SandboxLua)

func WithFuncName(funcName string) Option {
	return func(c *SandboxLua) {
		c.funcName = funcName
	}
}

type SandboxLua struct {
	state  *lua.LState
	cancel context.CancelFunc
	closed bool

	funcName string
	args     []lua.LValue
	result   []interface{}
}

func NewSandboxLua(opt ...Option) *SandboxLua {
	sandbox := &SandboxLua{
		state: lua.NewState(),
	}
	for _, op := range opt {
		op(sandbox)
	}
	return sandbox
}

func (l *SandboxLua) ParseArgs(args ...interface{}) error {
	nArgs := make([]lua.LValue, 0)
	for _, arg := range args {
		switch ar := arg.(type) {
		case float64:
			nArgs = append(nArgs, lua.LNumber(ar))
		case string:
			nArgs = append(nArgs, lua.LString(ar))
		case bool:
			nArgs = append(nArgs, lua.LBool(ar))
		case int:
			f := float64(ar)
			nArgs = append(nArgs, lua.LNumber(f))
		case int32:
			f := float64(ar)
			nArgs = append(nArgs, lua.LNumber(f))
		case int64:
			f := float64(ar)
			nArgs = append(nArgs, lua.LNumber(f))
		case float32:
			f := float64(ar)
			nArgs = append(nArgs, lua.LNumber(f))
		default:
			return TypeError
		}
	}
	l.args = nArgs
	return nil
}

func (l *SandboxLua) RunScript(code string) error {
	if l.closed {
		return CloseError
	}
	defer l.Close()
	if err := l.state.DoString(code); err != nil {
		return err
	}
	//非函数，直接执行
	if len(l.funcName) == 0 {
		return nil
	}
	//找有多数个返回值
	resultNum := utils.CountBackParameters(code)
	//需要函数范围才执行下面的内容
	if resultNum < 1 {
		return nil
	}
	// 调用 Lua 函数并获取返回值
	if err := l.state.CallByParam(lua.P{
		Fn:      l.state.GetGlobal(l.funcName),
		NRet:    resultNum, // 期望返回值个数
		Protect: true,
	}, l.args...); err != nil {
		return err
	}

	result := make([]interface{}, 0)
	for i := 0; i < resultNum; i++ {
		rept := l.state.Get(-1)
		value, err := toGoType(rept)
		if err != nil {
			return err
		}
		result = append(result, value)
		l.state.Pop(1)
	}

	l.result = reverse(result)
	return nil
}

func (l *SandboxLua) Timeout(timeout time.Duration) error {
	if l.closed {
		return CloseError
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	l.state.SetContext(ctx)
	l.cancel = cancel
	return nil
}

func reverse(slice []interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func (l *SandboxLua) Close() {
	if !l.closed {
		if l.cancel != nil {
			l.cancel()
		}
		l.state.Close()
		l.closed = true
	}
}

func (l *SandboxLua) GetResult() interface{} {
	return l.result
}

func toGoType(v lua.LValue) (interface{}, error) {
	if b, ok := v.(lua.LString); ok {
		return string(b), nil
	}
	if n, ok := v.(lua.LNumber); ok {
		return float64(n), nil
	}
	if n, ok := v.(lua.LBool); ok {
		return bool(n), nil
	}
	//if tbl, ok := v.(*lua.LTable); ok {
	//
	//}
	return nil, fmt.Errorf("cannot cast lua value to go value: %v", v)
}

// 按照 lua 里 if 语句的行为将任意 value 转换成 bool
func IsTrue(v lua.LValue) bool {
	switch v.Type() {
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTNil:
		return false
	default:
		return true
	}
}
