package chatapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func callGemini(prompt, apiKey string) (string, error) {
	// Prepare request body
	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	urlworking := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	fmt.Println("Using URL:", urlworking)

	req, err := http.NewRequest("POST", urlworking, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "No response generated", nil
}

func Runner() {
	// Get API key from environment variable (SECURE)
	apiKey := "AIzaSyCIuK9FepvqM6WfnCtP6fMFwjcctc_tVCg"

	fmt.Print("Enter your prompt: ")
	reader := bufio.NewReader(os.Stdin)
	userPrompt, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading input:", err)
	}
	userPrompt = strings.TrimSpace(userPrompt)

	if userPrompt == "" {
		log.Fatal("Prompt cannot be empty")
	}

	// Call Gemini API
	response, err := callGemini(userPrompt, apiKey)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n--- Gemini Response ---")
	fmt.Println(response)
}
