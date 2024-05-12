package main

import (
	"fmt"
	"com.github.zideajang.tinychain/prompt"
)

func main() {
	// Create a new prompt template for a joke about a specified content
	promptTemplate := prompt.FromTemplate("Tell me a %s joke about %s.")
	promptOne := promptTemplate.Format("funny", "chickens")
	fmt.Println(promptOne)

	// Create a simple prompt template for a generic joke
	promptTemplate = prompt.FromTemplate("Tell me a joke")
	promptTwo := promptTemplate.Format() // No variables needed here
	fmt.Println(promptTwo)
}