package turing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const apiURL = "http://127.0.0.1:11434/api/generate" // updated endpoint

var (
	rateLimiterMu   sync.Mutex
	lastRequestTime = time.Now().Add(-10 * time.Second)
)

// GetResponse : get response from ollama server
func GetResponse(input, modelName string) string {
	// Rate limiting to 1 request per 10 seconds
	rateLimiterMu.Lock()
	elapsed := time.Since(lastRequestTime)
	if elapsed < 10*time.Second {
		time.Sleep(10*time.Second - elapsed)
	}
	lastRequestTime = time.Now()
	rateLimiterMu.Unlock()

	data := []byte(fmt.Sprintf(`{"model": %s, "prompt": "%s", "stream": false}`, modelName, input))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "Could not reach Ollama server"
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return "Could not reach Ollama server"
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	responseRaw, _ := io.ReadAll(resp.Body)

	// Define a struct to parse the full JSON response from the API
	var res struct {
		Model              string `json:"model"`
		CreatedAt          string `json:"created_at"`
		Response           string `json:"response"`
		Done               bool   `json:"done"`
		DoneReason         string `json:"done_reason"`
		Context            []int  `json:"context"`
		TotalDuration      int64  `json:"total_duration"`
		LoadDuration       int64  `json:"load_duration"`
		PromptEvalCount    int    `json:"prompt_eval_count"`
		PromptEvalDuration int64  `json:"prompt_eval_duration"`
		EvalCount          int    `json:"eval_count"`
		EvalDuration       int64  `json:"eval_duration"`
	}

	err = json.Unmarshal(responseRaw, &res)
	if err != nil {
		log.Printf("Failed to parse response '%s': %v", responseRaw, err)
		// Return valid JSON so the caller can parse it
		return `{"response": "", "done": true}`
	}

	// Strip the <think></think> block from the response
	cleanResponse := strings.ReplaceAll(res.Response, "<think>", "")
	cleanResponse = strings.ReplaceAll(cleanResponse, "</think>", "")
	return cleanResponse
}
