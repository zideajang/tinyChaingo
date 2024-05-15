package tool

import (
    "fmt"
    "reflect"
)

type Tool struct {
    Name			string
    Description		string
    Args			map[string]map[String]string
    ReturnDirect 	bool
	Function		reflect.value
}

type ToolRunner interface {
    Run(args map[string]string) (interface{}, error)
}

func (t *Tool) Run(args map[string]string) (interface{}, error) {

    if t.Function.Kind() != reflect.Func {
        return nil, fmt.Errorf("provided function is not valid")
    }

    funcType := t.Function.Type()

    callArgs := make([]reflect.Value, funcType.NumIn())

    for i := 0; i < funcType.NumIn(); i++ {
        argName := funcType.In(i).Name()  
        if argValue, ok := args[argName]; ok {
            callArgs[i] = reflect.ValueOf(argValue)
        } else {
            return nil, fmt.Errorf("missing argument for %s", argName)
        }
    }

    results := t.Function.Call(callArgs)
    if len(results) != 1 {
        return nil, fmt.Errorf("function does not return exactly one result")
    }

    return results[0].Interface(), nil
}


func SearchFunc(prompt string) string {
	return "Search with " + prompt
}

