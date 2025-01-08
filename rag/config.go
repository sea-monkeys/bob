package rag

import (
	"encoding/json"
	"os"
)

type RagConfig struct {
	ChunkSize           int     `json:"chunkSize"`
	ChunkOverlap        int     `json:"chunkOverlap"`
	SimilarityThreshold float64 `json:"similarityThreshold"`
	MaxSimilarity       int     `json:"maxSimilarity"`
}

func LoadRagConfig(path string) (RagConfig, error) {
	// Load the json rag config file
	ragConfigFile, errRagConf := os.ReadFile(path)
	if errRagConf != nil {
		//log.Fatalf("😡 Error reading rag.json file: %v", errRagConf)
		return RagConfig{}, errRagConf
	}
	var ragConfig RagConfig
	errJsonRagConf := json.Unmarshal(ragConfigFile, &ragConfig)
	if errJsonRagConf != nil {
		//log.Fatalf("😡 Error unmarshalling rag.json file: %v", errJsonRagConf)
		return RagConfig{}, errJsonRagConf
	}
	return ragConfig, nil
}
