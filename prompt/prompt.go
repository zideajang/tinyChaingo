package prompt

import(
	"fmt"
	"regexp"
	"strings"
	"encoding/json"
)
import "tinychain/runnable"
	
	

type  IFormat interface{
	Format(InputVariables map[string]string) string
}


// String
type StringPromptTemplate struct{
	Template string
}

type PromptTemplate struct{
	StringPromptTemplate
	Next runnable.Runnable
	InputVariables []string
}

func (spt *StringPromptTemplate) Format(InputVariables map[string]string) string{
	result := spt.Template
	for k,v := range InputVariables{
		placeholder := "{" + k + "}"
		result = strings.ReplaceAll(result,placeholder,v)
	}
	return result
}

func (pt *PromptTemplate) SetNext(runnable runnable.Runnable) runnable.Runnable {
	pt.Next = runnable
	return runnable
}

func (pt *PromptTemplate) Invoke(input string) bool{
	var inputVariables map[string]string

	err := json.Unmarshal([]byte(input), &inputVariables)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false
	}

	if pt.Next != nil {
		return pt.Next.Invoke(pt.Format(inputVariables))
	}

	return false
	
	
}

func (pt *PromptTemplate) FromTemplate(template string){
	pt.Template = template
	variableRegex := regexp.MustCompile(`{([^]+)}`)
	matches := variableRegex.FindAllStringSubmatch(template,-1)
	for _, match := range matches{
		pt.InputVariables = append(pt.InputVariables,match[1])
	}
}



