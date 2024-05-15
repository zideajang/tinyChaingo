package outputparser

import "tinychain/runnable"

type OutputParser struct{
	Input string
}

func (pt *OutputParser) SetNext(runnable runnable.Runnable) runnable.Runnable {
	pt.Next = runnable
	return runnable
}


func (pt *PromptTemplate) Invoke(input string) bool{

	if pt.Next != nil {
		return pt.Next.Invoke(input)
	}

	return false
}