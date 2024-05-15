package main

import (
	"fmt"
	"reflect"
)


type User struct{
	Name string
	Age int
}

func main (){
	var x float64 = 3.1415926
	var tony User = User{"tony",28}
	_ = tony
	fmt.Println("type:",reflect.TypeOf(x))
	fmt.Println("value:", reflect.ValueOf(x))
}