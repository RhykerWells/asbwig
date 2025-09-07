package web

import (
	"encoding/json"
	"html/template"
	"strings"
)

var (
	templateFunctions = map[string]interface{}{
		"dict": dict,
		"seq": seq,
		"inSlice": inSlice,
		"toJson": toJson,
		"lower": lower,
	}
)

func dict(pairs ...interface{}) map[int]interface{} {
	result := make(map[int]interface{})
	for i := 0; i < len(pairs); i += 2 {
		key, _ := pairs[i].(int)
		result[key] = pairs[i+1]
	}
	return result
}


func seq(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func inSlice(val string, slice interface{}) bool {
	switch s := slice.(type) {
	case []string:
		for _, item := range s {
			if item == val {
				return true
			}
		}
	case []interface{}:
		for _, item := range s {
			if str, ok := item.(string); ok && str == val {
				return true
			}
		}
	}
	return false
}

func toJson(v interface{}) template.JS {
	b, _ := json.Marshal(v)
	return template.JS(b)
}

func lower(str string) string {
	return strings.ToLower(str)
}