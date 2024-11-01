package lua

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLuaRun(t *testing.T) {
	RunLua()
}

func TestK(t *testing.T) {
	code :=
		`function calculate(a, b)
    		local sum = a + b
    		local difference = a - b
			return sum, difference
		end`
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
	//loadString, err := L.LoadString(code)
	//if err != nil {
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

var code = `function calculate(a, b)
    		local sum = a + b
    		local difference = a - b
			return sum, difference
		end`

func TestWWWds(t *testing.T) {
	//WithFuncName("calculate"), WithResultNum(2)
	sandboxLua := NewSandboxLua()
	if err := sandboxLua.ParseArgs(8, 9); err != nil {
		panic(err)
	}
	if err := sandboxLua.RunScript(`print("hello")`); err != nil {
		panic(err)
	}
	result := sandboxLua.GetResult()
	fmt.Printf("Result: %v\n", result)
}
