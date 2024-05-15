package prompt

import (
	"fmt"
	"testing"
)

// 测试 StringPromptTemplate 的 Format 方法
func TestStringPromptTemplate_Format(t *testing.T) {
	stringTemplate := StringPromptTemplate{
		Template: "Tell me a {adjective} joke about {content}.",
	}

	inputVariables := map[string]string{
		"adjective": "funny",
		"content":   "chickens",
	}

	// 使用 stringTemplate 的 Format 方法
	result := stringTemplate.Format(inputVariables)

	// 预期结果字符串
	expected := "Tell me a funny joke about chickens."

	// 检查结果是否如预期
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	} else {
		fmt.Println("Test passed:", result)
	}
}
