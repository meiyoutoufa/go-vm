package python

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/meiyoutoufa/go-vm/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

func II() {
	getwd, _ := os.Getwd()
	fmt.Println(getwd)
	cmd := exec.Command("python", "./main.py")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	fmt.Println(string(output))
}

//go:embed template.tpl
var tmpl string

type SandboxPython struct {
	pythonVersion string
	funcName      string
	args          []string
	result        string
}

type tpl struct {
	Code      string
	CodeCited string
}

type Option func(c *SandboxPython)

func WithFuncName(funcName string) Option {
	return func(c *SandboxPython) {
		c.funcName = funcName
	}
}

func WithPythonVersion(pythonVersion string) Option {
	return func(c *SandboxPython) {
		c.pythonVersion = pythonVersion
	}
}

func NewSandboxPython(opts ...Option) *SandboxPython {
	sandbox := &SandboxPython{}
	for _, opt := range opts {
		opt(sandbox)
	}
	return sandbox
}

func (p *SandboxPython) ParseArgs(args ...interface{}) error {
	nargs := make([]string, 0)
	for _, arg := range args {
		argStr, err := utils.ConvertToString(arg)
		if err != nil {
			return err
		}
		nargs = append(nargs, argStr)
	}

	p.args = nargs
	return nil
}

func (p *SandboxPython) RunScript(code string) error {
	pyFileName, err := templatePythonFile(code, p.funcName, p.args...)
	if err != nil {
		return err
	}
	defer os.Remove(pyFileName)
	cmdArgs := make([]string, 0)
	cmdArgs = append(cmdArgs, "./main.py")
	cmdArgs = append(cmdArgs, p.args...)
	cmd := exec.Command(p.pythonVersion, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			// 捕获并打印标准错误输出
			fmt.Printf("Command failed with error: +%v", string(exitError.Stderr))
			return exitError
		}
		return err
	}
	fmt.Println(string(output))
	outPutStr := string(output)
	//"return " 后面的内容,没有说明不需要返回
	if utils.CountBackParameters(code) < 1 {
		return nil
	}

	space := strings.TrimSpace(outPutStr)
	split := strings.Split(space, "\n")
	fmt.Println(space)
	var result string
	if len(split) < 1 {
		result = split[0]
	} else {
		result = split[len(split)-1]
	}
	p.result = result
	return nil
}

func (p *SandboxPython) GetResult() interface{} {
	if p.result == "" {
		return nil
	}
	// 转换字符串为 float64
	value, err := strconv.ParseFloat(p.result, 64)
	if err != nil {
		return p.result
	}
	return value
}

func templatePythonFile(code, funcName string, args ...string) (string, error) {
	var tplObj tpl
	//没有funcName 则直接运行的内容
	if len(funcName) == 0 {
		tplObj = tpl{CodeCited: code}
	} else {
		parseFunc, err := getParseFunc(funcName, args...)
		if err != nil {
			return "", err
		}
		tplObj = tpl{Code: code, CodeCited: parseFunc}
	}

	t, err := template.New("tpl1").Parse(tmpl)
	if err != nil {
		return "", err
	}
	file, err := os.Create("./main.py")
	if err != nil {
		return "", err
	}
	defer file.Close()
	err = t.Execute(file, tplObj)
	return file.Name(), err
}

func getParseFunc(funcName string, args ...string) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("(")

	for i, argStr := range args {
		if i == len(args)-1 {
			sb.WriteString(fmt.Sprintf("%s", argStr))
			break
		}
		sb.WriteString(fmt.Sprintf("%s", argStr))
		sb.WriteString(",")
	}

	sb.WriteString(")")

	return funcName + sb.String(), nil
}

func CheckPythonInstalled() (string, error) {
	// 尝试使用 python3
	cmd := exec.Command("python3", "--version")
	if output, err := cmd.CombinedOutput(); err == nil {
		fmt.Println(strings.TrimSpace(string(output)))
		return "python3", nil
	}

	// 尝试使用 python
	cmd = exec.Command("python", "--version")
	if output, err := cmd.CombinedOutput(); err == nil {
		fmt.Println(strings.TrimSpace(string(output)))
		return "python", nil
	}

	return "", fmt.Errorf("Python is not installed")
}
