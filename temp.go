package main

type BasePromptTemplate struct {
	template string
	inputVariables []string
}

func (tmpl *BasePromptTemplate) Format(vars map[string]string) string {

}


	
import "fmt"
func main() {
    fmt.Println("hello world")
}