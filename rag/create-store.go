package rag

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/asellus"
	"github.com/sea-monkeys/bob/config"
	"github.com/sea-monkeys/daphnia"
)

func CreateVectorStore(ctx context.Context, config config.Config, ollamaClient *api.Client, ollamaRawUrl, embeddingsModel string) error {

	// Load the json rag config file
	ragConfig, errRagConf := LoadRagConfig(config.SettingsPath + "/rag.json")
	if errRagConf != nil {
		log.Fatalf("😡 Error loading rag.json file: %v", errRagConf)
		return errRagConf
	}

	// Initialize the vector store
	vectorStore := daphnia.VectorStore{}
	vectorStore.Initialize(config.SettingsPath + "/chunks.gob")

	// Read the content of the documents directory
	fmt.Println("📝🤖 using:", ollamaRawUrl, embeddingsModel, "for RAG.")
	fmt.Println("📝🤖 RAG Vector store creation in progress.")

	// Iterate over all the files in the content directory
	// and create embeddings for each file
	asellus.ForEveryFile(config.RagDocumentsPath, func(documentPath string) error {
		fmt.Println("📝 Creating embedding from document ", documentPath)

		// Read the content of the file
		document, err := asellus.ReadTextFile(documentPath)
		if err != nil {
			fmt.Println("😡 Error reading the content of the file:", err)
			//return err
		}
		//chunks := asellus.ChunkText(document, 2048, 512)
		// the values are defined in the ./bob/rag.json file
		chunks := asellus.ChunkText(document, ragConfig.ChunkSize, ragConfig.ChunkOverlap)

		fmt.Println("👋 Found", len(chunks), "chunks")

		// Create embeddings from documents and save them in the store
		for idx, chunk := range chunks {
			fmt.Println("📝 Creating embedding nb:", idx)
			fmt.Println("📝 Chunk:", chunk)

			req := &api.EmbeddingRequest{
				Model:  embeddingsModel,
				Prompt: chunk,
			}
			resp, errEmb := ollamaClient.Embeddings(ctx, req)
			if errEmb != nil {
				fmt.Println("😡 Error when calculating the embeddings:", errEmb)
				//return errEmb
			}

			// Save the embedding in the vector store
			_, err := vectorStore.Save(daphnia.VectorRecord{
				Prompt:    chunk,
				Embedding: resp.Embedding,
				Id:        documentPath + "-" + strconv.Itoa(idx),
				// The Id must be unique
			})

			//fmt.Println("📝 Embedding:", record.Embedding)

			if err != nil {
				fmt.Println("😡 Error when saving the embeddings:", err)
				//return err
			}
		}

		return nil
	})
	fmt.Println("📝🤖 RAG Vector store creation done 🎉.")
	return nil

}
