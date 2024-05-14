package runnable

type Runnable interface{
	Invoke(context string) string
}

type Node struct{
	RunnableNode Runnable
	Next *Node
}

type RunnableManager struct{
	Head *Node
	Tail *Node
	Context string
}

func (rm *RunnableManager) Add(runnable Runnable) {
    newNode := &Node{RunnableNode: runnable}
    if rm.Head == nil {
        rm.Head = newNode
        rm.Tail = newNode
    } else {
        rm.Tail.Next = newNode
        rm.Tail = newNode
    }
}

func (rm *RunnableManager) Run() {
    current := rm.Head
    for current != nil {
        rm.Context = current.RunnableNode.Invoke(rm.Context)
        current = current.Next
    }
}
