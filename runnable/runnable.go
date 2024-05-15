package runnable

// type RunnableContext struct {
//     Content string
// }

type Runnable interface{
    SetNext(runnable Runnable) Runnable
	Invoke(input string) bool
}


