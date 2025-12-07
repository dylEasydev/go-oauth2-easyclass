package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func ReadJSON[T any](fileName string) (*T, error) {
	baseDir, _ := os.Getwd()
	fullPath := path.Join(baseDir, "ressource/", fileName+".json")
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture fichier: %w", err)
	}
	defer file.Close()

	var data T
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("erreur decodage JSON: %w", err)
	}

	return &data, nil
}
