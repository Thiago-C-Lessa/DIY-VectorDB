package embedding

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt []string `json:"input"`
}

type OllamaEmbeddingResponse struct {
	Model     string         `json:"model"`
	Embedding [][768]float32 `json:"embeddings"`
}

func getEmbedding(model string, prompt []string) ([768]float32, error) {
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return [768]float32{}, err
	}

	resp, err := http.Post("http://ollama:11434/api/embed", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return [768]float32{}, err
	}
	defer resp.Body.Close()

	var ollamaResponse OllamaEmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&ollamaResponse)
	if err != nil {
		return [768]float32{}, err
	}

	return ollamaResponse.Embedding[0], nil
}

func GenereteEmbedding(key string) ([768]float32, error) {
	model := "embeddinggemma:300m"
	prompt := key
	var input []string
	input = append(input, prompt)

	emb, err := getEmbedding(model, input)
	if err != nil {
		return [768]float32{}, err
	}

	return emb, nil
}
