package python

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	parseFunc, err := getParseFunc("add", "ss", "5", "4")
	fmt.Println(parseFunc)
	fmt.Println(err)
}
