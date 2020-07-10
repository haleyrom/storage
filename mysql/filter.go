package mysql

// SQL注入过滤

import (
	"fmt"
	"github.com/haleyrom/storage/logging"

	"strings"
)

//模糊匹配列表
var filter = []string{"--", ";", "'", "(", ")", "="}

//完整匹配列表
var filterFull = []string{"information_schema", "sleep", "drop", "truncate", "delete", "insert", "update", "rlike"}

// SqCondSprintf SqCondSprintf
func SqCondSprintf(format string, a ...interface{}) string {
	for _, param := range a {
		switch param.(type) {
		case string:
			if !filterCheck(param.(string)) {
				logging.Error(nil, "sql attack:%s", param.(string))
				panic("SqCondSprintf")
				return ""
			}
		}
	}
	return fmt.Sprintf(format, a...)
}

// filterCheck filterCheck
func filterCheck(val string) bool {
	//以空格拆分字符串
	vs := strings.Split(val, " ")
	for _, v := range vs {
		//符号类模糊匹配
		for _, f := range filter {
			if strings.Index(strings.ToLower(v), f) != -1 {
				return false
			}
		}

		//字符串类完全匹配
		for _, f := range filterFull {
			if f == strings.ToLower(v) {
				return false
			}
		}
	}

	return true
}
