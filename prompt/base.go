package prompt

import (
	"fmt"
	"strings"
)

type StringPromtTemplate struct{
	template string
	inputVariables []string
}

func NewStringPromptTemplate(template string, inputVariables []string) *StringPromptTemplate {
	return &StringPromptTemplate{
		template:       template,
		inputVariables: inputVariables,
	}
}

func (spt *StringPromptTemplate) Format(vars map[string]string) string {
	t := spt.template
	for key, val := range vars {
		placeholder := "{" + key + "}"
		t = strings.ReplaceAll(t, placeholder, val)
	}
	return t
}

type PromptTemplate struct {
	template string
}

func FromTemplate(template string) *PromptTemplate {
	return &PromptTemplate{template: template}
}

func (pt *PromptTemplate) Format(vars ...interface{}) string {
	return fmt.Sprintf(pt.template, vars...)
}