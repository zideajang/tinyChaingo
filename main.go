package main

import (
	//"fmt"
	// "tinychain/prompt"
	"tinychain/prompt"
	//"tinychain/runnable"
	// "tinychain/llm"
)

func main(){

	promptTemplate := &prompt.PromptTemplate{}

	promptTemplate.Template = "Tell me a {adjective} joke about {content}."


	//promptTemplate
	

	// ollama := &llm.OllamaChatModel{}

	// prompt := `[{"role": "user","content": "why is the sky blue?"}]`

	// response := ollama.Invoke(prompt)

	// fmt.Println("Response from Ollama Chat Model:", response)

	// stringTemplate := prompt.StringPromptTemplate{
	// 	Template: "Tell me a {adjective} joke about {content}.",
	// }
	
	
	// inputVariables := map[string]string{
	// 	"adjective":"funny",
	// 	"content":"chickens",
	// }
	
	// promptTemplate := prompt.PromptTemplate{}

	// _ = promptTemplate

	// result := stringTemplate.Format(inputVariables)
	// fmt.Println(result)
}
