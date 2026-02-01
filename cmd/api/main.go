package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ollamaRequest struct {
	Model  string   `json:"model"`
	Prompt []string `json:"input"`
}

type OllamaEmbeddingResponse struct {
	Model     string      `json:"model"`
	Embedding [][]float64 `json:"embeddings"`
}

func GetEmbedding(model string, prompt []string) error {
	reqBody := ollamaRequest{
		Model:  model,
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://ollama:11434/api/embed", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ollamaResponse OllamaEmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&ollamaResponse)
	if err != nil {
		return err
	}

	fmt.Println(len(ollamaResponse.Embedding[0]))
	fmt.Println(ollamaResponse.Embedding)

	return nil
}

func main() {
	model := "embeddinggemma:300m"
	prompt := "Texto para gerar embedding"
	var input []string
	input = append(input, prompt)

	err := GetEmbedding(model, input)
	if err != nil {
		fmt.Println(err)
	}
}
