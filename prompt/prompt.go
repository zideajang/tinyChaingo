package prompt

import(
	"regexp"
	"strings"
)

type  IFormat interface{
	Format(InputVariables map[string]string) string
}


// String
type StringPromptTemplate struct{
	Template string
}

type PromptTemplate struct{
	StringPromptTemplate
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

func (pt *PromptTemplate) Invoke(input string){
	pt.FromTemplate(input)
}

func (pt *PromptTemplate) FromTemplate(template string){
	pt.Template = template
	variableRegex := regexp.MustCompile(`{([^]+)}`)
	matches := variableRegex.FindAllStringSubmatch(template,-1)
	for _, match := range matches{
		pt.InputVariables = append(pt.InputVariables,match[1])
	}
}



