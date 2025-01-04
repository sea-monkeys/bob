package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/asellus"
	"github.com/sea-monkeys/daphnia"
)

// TODO: check if the model is loaded / exists
// TODO: add a waiting message
// TODO: add an option for the conversational memory
// TODO: add RAG features: https://k33g.hashnode.dev/rag-from-scratch-with-go-and-ollama?source=more_series_bottom_blogs
// TODO: generate the report and its content at the same time (streaming)
// TODO: add a command to generate project files

type RagConfig struct {
	ChunkSize           int     `json:"chunkSize"`
	ChunkOverlap        int     `json:"chunkOverlap"`
	SimilarityThreshold float64 `json:"similarityThreshold"`
	MaxSimilarity       int     `json:"maxSimilarity"`
}

type Config struct {
	PromptPath          string
	ContextPath         string // for this one check if the file exists
	ToolsInvocationPath string
	JsonSchemaPath      string

	SettingsPath     string
	OutputPath       string
	RagDocumentsPath string // for RAG

	// to override the system and user questions
	System string
	User   string
}

func validatePaths(config Config) error {
	// Check if prompt file exists
	if _, err := os.Stat(config.PromptPath); err != nil {
		return fmt.Errorf("prompt file not found: %v", err)
	}

	// Check if settings directory exists
	if info, err := os.Stat(config.SettingsPath); err != nil {
		return fmt.Errorf("settings directory not found: %v", err)
	} else if !info.IsDir() {
		return fmt.Errorf("settings path must be a directory")
	}

	// Check if output directory exists
	outputDir := filepath.Dir(config.OutputPath)
	if info, err := os.Stat(outputDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		} else {
			return fmt.Errorf("error checking output directory: %v", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("output path parent must be a directory")
	}

	return nil
}

func loadRagConfig(path string) (RagConfig, error) {
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

var (
	FALSE = false
	TRUE  = true
)

//go:embed version.txt
var versionTxt []byte

func main() {
	config := Config{}

	// Define command line flags
	flag.StringVar(&config.PromptPath, "prompt", "prompt.md", "Path to prompt file")

	flag.StringVar(&config.SettingsPath, "settings", ".bob", "Path to settings directory")
	flag.StringVar(&config.OutputPath, "output", "report.md", "Path to output file")
	flag.StringVar(&config.RagDocumentsPath, "rag", "", "Path to content directory for RAG")

	flag.StringVar(&config.ToolsInvocationPath, "tools-invocation", "tools.invocation.md", "Path to tools invocation file")
	flag.StringVar(&config.JsonSchemaPath, "json-schema", "schema.json", "Path to JSON schema file")
	flag.StringVar(&config.ContextPath, "context", "context.md", "Path to context file")

	flag.StringVar(&config.System, "system", "", "System instructions")
	flag.StringVar(&config.User, "user", "", "User question")

	// Version flag
	version := flag.Bool("version", false, "Display version information")

	// use bob --tools to invoke tools
	toolsInvocation := flag.Bool("tools", false, "Tools invocation")
	// use bob --schema to use a JSON schema
	jsonSchema := flag.Bool("schema", false, "JSON schema")

	// Parse command line arguments
	flag.Parse()

	// Check for version flag
	if *version {
		fmt.Println(string(versionTxt))
		os.Exit(0)
	}

	// Validate required flags
	if config.PromptPath == "" || config.SettingsPath == "" || config.OutputPath == "" {
		fmt.Println("Usage: bob --prompt path_to_prompt_file --settings path_to_settings_directory --output path_to_output_file")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate paths
	if err := validatePaths(config); err != nil {
		fmt.Printf("😡 Error: %v\n", err)
		os.Exit(1)
	}

	/*
		fmt.Printf("Processing with:\n")
		fmt.Printf("  Prompt: %s\n", config.PromptPath)
		fmt.Printf("  Settings: %s\n", config.SettingsPath)
		fmt.Printf("  Output: %s\n", config.OutputPath)
	*/

	// Main logic
	ctx := context.Background()

	errEnv := godotenv.Load(config.SettingsPath + "/.env")
	if errEnv != nil {
		log.Fatalf("😡 Error loading .env file: %v", errEnv)
		// Fatalf is equivalent to [Printf] followed by a call to os.Exit(1).
	}

	var ollamaRawUrl string
	if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
		ollamaRawUrl = "http://localhost:11434"
	}

	var model string
	if model = os.Getenv("LLM"); model == "" {
		model = "qwen2.5:0.5b"
	}
	var toolsModel string
	if toolsModel = os.Getenv("TOOLS_LLM"); toolsModel == "" {
		toolsModel = "qwen2.5:0.5b"
	}

	var embeddingsModel string
	if embeddingsModel = os.Getenv("EMBEDDINGS_LLM"); embeddingsModel == "" {
		embeddingsModel = "snowflake-arctic-embed:33m"
	}

	url, _ := url.Parse(ollamaRawUrl)

	fmt.Println("📣🤖 using:", ollamaRawUrl, model, "for Chat completion")

	// Model settings
	// Configuration
	modelConfigFile, errConf := os.ReadFile(config.SettingsPath + "/settings.json")
	if errConf != nil {
		log.Fatalf("😡 Error reading settings.json file: %v", errConf)
	}

	var modelConfig map[string]interface{}
	errJsonConf := json.Unmarshal(modelConfigFile, &modelConfig)
	if errJsonConf != nil {
		log.Fatalf("😡 Error unmarshalling settings.json file: %v", errConf)
	}

	ollamaClient := api.NewClient(url, http.DefaultClient)

	// ==========================================================
	// RAG Creation of the Vector Store
	// ==========================================================
	// create the vector store in .bob
	// then Bob will be able to detect if he needs to use it
	// Run it: go run ../../main.go --rag ./content
	if config.RagDocumentsPath != "" {

		// Load the json rag config file
		ragConfig, errRagConf := loadRagConfig(config.SettingsPath + "/rag.json")
		if errRagConf != nil {
			log.Fatalf("😡 Error loading rag.json file: %v", errRagConf)
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
				fmt.Println("😡:", err)
				// TODO: handle error
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
					fmt.Println("😡:", errEmb)
					// TODO: handle error
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
					fmt.Println("😡:", err)
					// TODO: handle error

				}
			}

			return nil
		})
		fmt.Println("📝🤖 RAG Vector store creation done 🎉.")
		os.Exit(0)
	}

	var systemInstructions, userQuestion string

	if config.System != "" {
		systemInstructions = config.System
	} else {
		// Load the content of the instructions.md file
		instructions, errInstruct := os.ReadFile(config.SettingsPath + "/instructions.md")
		if errInstruct != nil {
			log.Fatalf("😡 Error reading instructions file: %v", errInstruct)
		}
		systemInstructions = string(instructions)
	}

	if config.User != "" {
		userQuestion = config.User
	} else {
		// Load the content of the prompt.md file
		prompt, errPrompt := os.ReadFile(config.PromptPath)
		if errPrompt != nil {
			log.Fatalf("😡 Error reading prompt file: %v", errPrompt)
		}
		userQuestion = string(prompt)
	}

	messages := []api.Message{}
	messages = append(messages, api.Message{Role: "system", Content: systemInstructions})

	// ==========================================================
	// Context
	// ==========================================================
	var contextContent []byte
	// Check if the context file exists
	if _, err := os.Stat(config.ContextPath); err == nil {
		// Load the content of the context.md file
		var errContext error
		contextContent, errContext = os.ReadFile(config.ContextPath)
		if errContext != nil {
			log.Fatalf("😡 Error reading context file: %v", errContext)
		}
		//fmt.Println("📝 Context:", string(contextContent))
	}
	if string(contextContent) != "" {
		messages = append(messages, api.Message{Role: "system", Content: string(contextContent)})
	}

	// ==========================================================
	// Tools
	// ==========================================================
	toolsContext := ""

	if *toolsInvocation {
		toolsContext += "<documents>"
		// Tool invocation
		fmt.Println("🛠️🤖 using:", ollamaRawUrl, toolsModel, "for tools")

		// Read tools
		toolsConfigFile, errToolsConf := os.ReadFile(config.SettingsPath + "/tools.json")
		if errToolsConf != nil {
			log.Fatalf("😡 Error reading tools.json file: %v", errToolsConf)
		}
		var toolsList api.Tools
		errJsonToolsConf := json.Unmarshal(toolsConfigFile, &toolsList)
		if errJsonToolsConf != nil {
			log.Fatalf("😡 Error unmarshalling tools.json file: %v", errJsonToolsConf)
		}

		// Load the content of the tools.invocation.md file
		toolsPrompt, errPrompt := os.ReadFile(config.ToolsInvocationPath)
		if errPrompt != nil {
			log.Fatalf("😡 Error reading tools.invocation file: %v", errPrompt)
		}
		tools := strings.Split(string(toolsPrompt), "---")
		//fmt.Println("🛠️", tools)

		// Tools Prompt construction
		messagesForTools := []api.Message{}
		for _, tool := range tools {
			messagesForTools = append(messagesForTools, api.Message{Role: "user", Content: tool})
		}

		req := &api.ChatRequest{
			Model:    toolsModel,
			Messages: messagesForTools,
			Options: map[string]interface{}{
				"temperature": 0.0,
			},
			Tools:  toolsList,
			Stream: &FALSE,
		}

		err := ollamaClient.Chat(ctx, req, func(resp api.ChatResponse) error {

			for _, toolCall := range resp.Message.ToolCalls {
				fmt.Println("🛠️", toolCall.Function.Name, toolCall.Function.Arguments)

				// Convert map to slice of arguments
				cmdArgs := []string{config.SettingsPath + "/" + toolCall.Function.Name + ".sh"}
				for _, v := range toolCall.Function.Arguments {
					cmdArgs = append(cmdArgs, v.(string))
				}

				cmd := exec.Command("bash", cmdArgs...)
				output, err := cmd.Output()
				if err != nil {
					panic(err)
				}
				//fmt.Println("🤖", string(output))

				// Add the output to the context
				toolsContext += "<document>" + string(output) + "</document>"

				//messages = append(messages, api.Message{Role: "system", Content: string(output)})

			}
			toolsContext += "</documents>"
			fmt.Println()

			//fmt.Println("🤖", promptContext)

			//messages = append(messages, api.Message{Role: "system", Content: "CONTEXT:\n" + promptContext})
			return nil
		})

		if err != nil {
			log.Fatalln("😡", err)
		}

	} // end of tool invocation
	// ==========================================================

	// Prompt construction
	if toolsContext != "" {
		messages = append(messages, api.Message{Role: "system", Content: toolsContext})
		//userQuestion = promptContext + "\n\n" + userQuestion
	}

	/*
		messages = []api.Message{
			{Role: "system", Content: systemInstructions},
			{Role: "user", Content: userQuestion},
		}
	*/

	var req *api.ChatRequest

	if *jsonSchema {
		messages = append(messages, api.Message{Role: "user", Content: userQuestion})

		// Read the content of the schema.json file
		schema, errSchema := os.ReadFile(config.JsonSchemaPath)
		if errSchema != nil {
			log.Fatalf("😡 Error reading schema file: %v", errSchema)
		}
		// TMP
		//fmt.Println("🤖 using:", schema)
		req = &api.ChatRequest{
			Model:    model,
			Messages: messages,
			Options:  modelConfig,
			Stream:   &FALSE,
			Format:   json.RawMessage(schema),
		}

	} else { // classic chat completion

		// ==========================================================
		// Check if we need to use the vector store
		// ==========================================================

		// check if chunks.gob exists
		_, err := os.Stat(config.SettingsPath + "/chunks.gob")
		if err == nil { // then time to load the vector store and search for the closest chunks
			
			fmt.Println("📝🤖 using:", ollamaRawUrl, embeddingsModel, "for RAG.")

			// Load the json rag config file
			ragConfig, errRagConf := loadRagConfig(config.SettingsPath + "/rag.json")
			if errRagConf != nil {
				log.Fatalf("😡 Error loading rag.json file: %v", errRagConf)
			}

			// Load the vector store
			vectorStore := daphnia.VectorStore{}
			vectorStore.Initialize(config.SettingsPath + "/chunks.gob")

			question := userQuestion
			// Embbeding of the question - search for the closest chunk(s)
			reqEmbedding := &api.EmbeddingRequest{
				Model:  embeddingsModel,
				Prompt: question,
			}
			resp, errEmb := ollamaClient.Embeddings(ctx, reqEmbedding)
			if errEmb != nil {
				fmt.Println("😡:", errEmb)
				// TODO: handle error
			}
			embeddingFromQuestion := daphnia.VectorRecord{
				Prompt:    question,
				Embedding: resp.Embedding,
			}

			// the values are defined in the ./bob/rag.json file
			//similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.75, 50)
			//similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 10)
			similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, ragConfig.SimilarityThreshold, ragConfig.MaxSimilarity)
			if errSim != nil {
				fmt.Println("😡:", errSim)
				// TODO: handle error
			}

			/*
				for _, similarity := range similarities {
					fmt.Println()
					fmt.Println("Cosine distance:", similarity.CosineSimilarity)
					fmt.Println(similarity.Prompt)
				}
			*/

			if len(similarities) == 0 {
				fmt.Println("😠 No similarities found")
			} else {
				fmt.Println("🎉 number of similarities:", len(similarities))
			}

			// === prepare the ragContext for answering question ===
			//merge similarities into a single string
			ragContext := ""
			for _, similarity := range similarities {
				ragContext += similarity.Prompt + " "
			}

			//fmt.Println("📝 Context:", ragContext)

			messages = append(messages, api.Message{Role: "system", Content: "CONTEXT:\n" + ragContext})

		} // end of similarites search

		messages = append(messages, api.Message{Role: "user", Content: userQuestion})

		//fmt.Println(messages)

		req = &api.ChatRequest{
			Model:    model,
			Messages: messages,
			Options:  modelConfig,
			Stream:   &TRUE,
		}
	}

	// Send the request to the server
	answer := ""
	errCompletion := ollamaClient.Chat(ctx, req, func(resp api.ChatResponse) error {
		answer += resp.Message.Content
		fmt.Print(resp.Message.Content)
		return nil
	})

	if errCompletion != nil {
		log.Fatalf("😡 Completion error: %v", errCompletion)
	}

	// generate a markdown file from the value of answer
	errOutput := os.WriteFile(config.OutputPath, []byte(answer), 0644)
	if errOutput != nil {
		log.Fatalf("😡 Error writing output file: %v", errOutput)
	}
	fmt.Println()
}
