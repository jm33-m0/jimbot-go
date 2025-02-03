package turing

import (
	"bytes"
	"encoding/json" // added import
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiURL = "http://127.0.0.1:11434/chat" // updated endpoint

// GetResponse : get response from ollama server
func GetResponse(input string) string {
	// Updated payload with prompt field for Ollama
	data := []byte(fmt.Sprintf(`{"prompt": "%s"}`, input))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Sprintf("Could not reach Ollama server: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Could not reach Ollama server: %v", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	responseRaw, _ := io.ReadAll(resp.Body)

	// Define a struct to parse the JSON response using stdlib
	var res struct {
		Response string `json:"response"`
	}

	err = json.Unmarshal(responseRaw, &res)
	if err != nil {
		return "Failed to parse response"
	}
	return res.Response
}
