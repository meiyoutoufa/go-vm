package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CountParameters 函数根据语言类型提取参数个数
func CountParameters(funcStr string, lang string) int {
	var re *regexp.Regexp

	switch lang {
	case "lua":
		re = regexp.MustCompile(`function \w+\(([^)]*)\)`)
	case "python":
		re = regexp.MustCompile(`def \w+\(([^)]*)\):`)
	case "javascript":
		re = regexp.MustCompile(`function \w*\(([^)]*)\)|\(\s*([^)]*)\s*\)\s*=>`)
	default:
		return 0
	}

	match := re.FindStringSubmatch(funcStr)
	if len(match) > 1 {
		// 处理函数声明参数
		params := strings.Split(match[1], ",")
		//for i := range params {
		//	params[i] = strings.TrimSpace(params[i]) // 去除空格
		//}
		return len(params) // 返回参数个数
	}

	return 0 // 如果没有找到参数，返回 0
}

func CountBackParameters(code string) int {
	re := regexp.MustCompile(`(?i)return (.*)`) // (?i) 使匹配不区分大小写
	match := re.FindStringSubmatch(code)
	// 检查是否找到匹配
	if len(match) > 1 {
		//fmt.Println("Found:", match[1]) // match[1] 是 "return " 后面的内容
		split := strings.Split(match[1], ",") //使用,逗号识别返回个数
		return len(split)
	}
	return 0
}

func ConvertToString(i interface{}) (string, error) {
	switch v := i.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case float64, float32:
		return fmt.Sprintf("%f", v), nil
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", errors.New("args convert type error")
	}
}
