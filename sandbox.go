package go_vm

import (
	"fmt"
	"github.com/meiyoutoufa/go-vm/javascript"
	"github.com/meiyoutoufa/go-vm/lua"
	"github.com/meiyoutoufa/go-vm/python"
	"github.com/meiyoutoufa/go-vm/utils"
)

type sandbox interface {
	ParseArgs(args ...interface{}) error
	RunScript(code string) error
	GetResult() interface{}
}

const (
	Lua        = "lua"
	Javascript = "javascript"
	Python     = "python"
)

func RunScript(lang string, funcName, code string, args ...interface{}) (interface{}, error) {
	if err := checkArgsNumber(lang, code, args...); err != nil {
		return nil, err
	}
	box, err := getSandbox(lang, funcName)
	if err != nil {
		return nil, err
	}
	if err = box.ParseArgs(args...); err != nil {
		return nil, err
	}
	if err = box.RunScript(code); err != nil {
		return nil, err
	}
	result := box.GetResult()

	return result, nil
}

func getSandbox(lang string, funcName string) (sandbox, error) {
	if lang == Lua {
		var opts []lua.Option
		if len(funcName) != 0 {
			opts = append(opts, lua.WithFuncName(funcName))
		}
		return lua.NewSandboxLua(opts...), nil
	}

	if lang == Javascript {
		var opts []javascript.Option
		if len(funcName) != 0 {
			opts = append(opts, javascript.WithFuncName(funcName))
		}
		return javascript.NewSandboxJs(opts...), nil
	}

	if lang == Python {
		var opts []python.Option
		if len(funcName) != 0 {
			opts = append(opts, python.WithFuncName(funcName))
		}
		pythonVersion, err := python.CheckPythonInstalled()
		if err != nil {
			return nil, err
		}
		opts = append(opts, python.WithPythonVersion(pythonVersion))

		return python.NewSandboxPython(opts...), nil
	}
	return nil, fmt.Errorf("lang error")
}

func checkArgsNumber(lang string, code string, args ...interface{}) error {
	parameters := utils.CountParameters(code, lang)
	if parameters != len(args) {
		return fmt.Errorf("wrong number of parameters")
	}
	return nil
}
