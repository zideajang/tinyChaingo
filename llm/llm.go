package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


// requst body -> OllamachatModel -> response

/**
curl http://localhost:11434/api/chat -d '{
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    }
  ],
  "stream": false
}'

*/

// define model 
type OllamaChatModel struct{}

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model string `json:"model"`
	Messages  []Message `json:"messages"`
	Stream bool `json:"stream"`
}

type ResponseBody struct {
	Message struct {
		Content string `json:"content"`
	} `json:message`
}

func (ocm *OllamaChatModel) Invoke(prompt string) string {

	fmt.Println(prompt)
	var messages []Message
	if err := json.Unmarshal([]byte(prompt),&messages); err != nil{
		fmt.Println("Error parsing prompt:",err)
		return ""
	}

	requestBody := RequestBody{
		Model: "llama3",
		Messages: messages,
		Stream: false,
	}

	fmt.Println(requestBody)

	reqBytes, err := json.Marshal(requestBody)
	if err != nil{
		fmt.Println("Error encoding request:",err )
		return ""
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(reqBytes))
    if err != nil {
		
		fmt.Println("Http Post response:", err)
		return ""
	}
	
	if err != nil {
		fmt.Println("Error sending request:", err)
        return ""
    }

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error sending ReadAll:", err)
        return ""
    }
	
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
        fmt.Println("Error parsing response:", err)
        return ""
    }

	fmt.Println(responseBody)
	return responseBody.Message.Content

}

