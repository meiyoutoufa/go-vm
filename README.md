### go exec script tool

如果使用python，需要安装python解释器，并且需要加入到path中,同时需要严格按照python要求的缩进的代码格式

Import a package.
```shell
import (
    "github.com/meiyoutoufa/go-vm"
)
```
install 
```shell
go get github.com/meiyoutoufa/go-vm
```




example

```go

import (
	"fmt"
	"testing"
)

func TestLuaNotFunction(t *testing.T) {
	script, err := RunScript(Lua, "", `print("hello, world")`)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(script)
}

func TestJavaScriptNotFunction(t *testing.T) {
	code := ` console.log('hello, world')`
	script, err := RunScript(Javascript, "", code)

	if err != nil {
		panic(err)
	}
	fmt.Println(script)
}

func TestPythonNotFunction(t *testing.T) {
	code := `print("hello, world")`
	script, err := RunScript(Python, "", code)

	if err != nil {
		panic(err)
	}
	fmt.Println(script)
}

func TestLuaFunction(t *testing.T) {
	var code = `function calculate(a, b)
    		local sum = a + b
    		local difference = a - b
			return sum, difference
		end`
	script, err := RunScript(Lua, "calculate", code, 1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(script)
}

func TestJSFunction(t *testing.T) {
	var code = `function calculate(a, b) {
			sum = a + b;
            return sum;
        }`
	script, err := RunScript(Javascript, "calculate", code, 5, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(script)
}

func TestPythonFunction(t *testing.T) {
	code :=
		`def calculate(a,b):
    print(a+b)
    return a+b`
	script, err := RunScript(Python, "calculate", code, 4, 5)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(script)
}
```

