package IRepositoryShare

import (
	"context"
	"net/http"
)

type IAI interface {
	GenerateBgeM3Embedding(text string) ([]float32, error)
	CallOllamaModel(ctx context.Context, prompt string, imageBase64 []string, httpClient *http.Client) (string, error)
}
