// Source code
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

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/bob/config"
	"github.com/sea-monkeys/bob/rag"
	"github.com/sea-monkeys/bob/tools"
	"github.com/sea-monkeys/bob/utilities"
	"github.com/sea-monkeys/daphnia"
)

// TODO: check if the model is loaded / exists
// TODO: add a waiting message
// TODO: add an option for the conversational memory
// TODO: generate the report and its content at the same time (streaming)
// TODO: add several files to the messages?


var (
	FALSE = false
	TRUE  = true
)

//go:embed version.txt
var versionTxt []byte

func main() {
	config := config.Config{}

	// Define command line flags
	flag.StringVar(&config.PromptPath, "prompt", "prompt.md", "Path to prompt file")

	flag.StringVar(&config.SettingsPath, "settings", ".bob", "Path to settings directory")
	flag.StringVar(&config.OutputPath, "output", "report.md", "Path to output file")
	flag.StringVar(&config.RagDocumentsPath, "rag", "", "Path to content directory for RAG")

	flag.StringVar(&config.ToolsInvocationPath, "tools-invocation", "tools.invocation.md", "Path to tools invocation file")
	flag.StringVar(&config.JsonSchemaPath, "json-schema", "schema.json", "Path to JSON schema file")
	flag.StringVar(&config.ContextPath, "context", "context.md", "Path to context file")

	flag.StringVar(&config.AddToMessages, "add-to-messages", "", "Add to messages")

	// Project structure
	flag.StringVar(&config.CreateProjectPathName, "create", "", "Project path name")
	flag.StringVar(&config.KindOfProject, "kind", "chat", "Kind of project")

	flag.StringVar(&config.System, "system", "", "System instructions")
	flag.StringVar(&config.User, "user", "", "User question")

	// Version flag
	version := flag.Bool("version", false, "Display version information")

	// use bob --tools to invoke tools
	toolsInvocation := flag.Bool("tools", false, "Tools invocation")
	// use bob --schema to use a JSON schema
	jsonSchema := flag.Bool("schema", false, "JSON schema")

	asUserMessage := flag.Bool("as-user", false, "As user message")
	asSystemMessage := flag.Bool("as-system", false, "As system message")

	beforeQuestion := flag.Bool("before-question", false, "Before user question")
	afterQuestion := flag.Bool("after-question", false, "After user question")
	// Parse command line arguments
	flag.Parse()

	// ==========================================================
	// 👷 Start of Project Creation: create project structure
	// ==========================================================
	/* Command examples:
	```bash
	bob --create demo
	````
	*/
	if config.CreateProjectPathName != "" { // Create a project structure and Exit
		err := utilities.CreateProject(config)
		if err != nil {
			os.Exit(1)
		} else {
			// Project created, exit
			os.Exit(0)
		}
	}
	// End of Project Creation

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
	// 👷 RAG Creation of the Vector Store
	// ==========================================================
	// create the vector store in .bob
	// then Bob will be able to detect if he needs to use it
	// Run it: go run ../../main.go --rag ./content
	/* Command examples:
	```bash
	bob --settings samples/chronicles-of-aethelgard/.bob \
	--rag samples/chronicles-of-aethelgard/content

	bob --rag ./content
	````
	*/
	if config.RagDocumentsPath != "" { // Create a vector store and Exit

		err := rag.CreateVectorStore(ctx, config, ollamaClient, ollamaRawUrl, embeddingsModel)
		if err != nil {
			os.Exit(1)
		} else {
			// Vector store created, exit
			os.Exit(0)
		}

	} // end of vector store creation

	// ==========================================================
	// 📝 Prepare the messages list for the completion
	// ==========================================================
	/* Command examples:
	```bash
	bob --system "You are an expert in Geography" --user "What is the capital of France?"
	```
	*/

	var systemInstructions, userQuestion string

	if config.System != "" { // override the system instructions contained in the instructions.md file
		systemInstructions = config.System
	} else {
		// Load the content of the instructions.md file
		instructions, errInstruct := os.ReadFile(config.SettingsPath + "/instructions.md")
		if errInstruct != nil {
			log.Fatalf("😡 Error reading instructions file: %v", errInstruct)
		}
		systemInstructions = string(instructions)
	}

	if config.User != "" { // override the user question contained in the prompt.md file
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
	// 📝 Add Context
	// ==========================================================
	// You cand add a context to the conversation
	/* Command examples:
	```bash
	bob --context /path/to/context.md
	# by default if a context.md file exists at the root of the project, it will be used
	```
	*/

	var contextContent []byte
	// Check if the context file exists in the defined path
	// You must define the path to the context file if you want to use it
	// Then the path could be different from the prompt.md file path
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
	// 🛠️ Tools
	// ==========================================================
	toolsContext := ""

	if *toolsInvocation { // bob --tools
		var err error
		toolsContext, err = tools.ToolsInvocation(ctx, config, ollamaClient, ollamaRawUrl, toolsModel)
		if err != nil {
			fmt.Println("😡 Error invoking tools:", err)
			os.Exit(1)
		}

	} // end of tool invocation


	var req *api.ChatRequest

	if *jsonSchema { // bob --schema
		messages = append(messages, api.Message{Role: "user", Content: userQuestion})

		// Read the content of the schema.json file
		schema, errSchema := os.ReadFile(config.JsonSchemaPath)
		if errSchema != nil {
			fmt.Println("😡 Error reading schema file:", errSchema)
			os.Exit(1)
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

	} else { // classic chat completion with RAG or not

		// ==========================================================
		// Check if we need to use the vector store
		// ==========================================================

		// check if chunks.gob exists
		_, err := os.Stat(config.SettingsPath + "/chunks.gob")
		if err == nil { // then time to load the vector store and search for the closest chunks

			fmt.Println("📝🤖 using:", ollamaRawUrl, embeddingsModel, "for RAG.")

			// Load the json rag config file
			ragConfig, errRagConf := rag.LoadRagConfig(config.SettingsPath + "/rag.json")
			if errRagConf != nil {
				fmt.Println("😡 Error loading rag.json file:", errRagConf)
				os.Exit(1)
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
				fmt.Println("😡 Error with embeddings request", errEmb)
				os.Exit(1)
			}
			embeddingFromQuestion := daphnia.VectorRecord{
				Prompt:    question,
				Embedding: resp.Embedding,
			}

			// The values are defined in the ./bob/rag.json file
			//similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.75, 50)
			//similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.3, 10)
			similarities, errSim := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, ragConfig.SimilarityThreshold, ragConfig.MaxSimilarity)
			if errSim != nil {
				fmt.Println("😡 Error when searching the similarities", errSim)
				os.Exit(1)
			}

			if len(similarities) == 0 {
				fmt.Println("😠 No similarities found")
			} else {
				fmt.Println("🎉 number of similarities:", len(similarities))
			}

			// === prepare the ragContext for answering question ===
			// merge similarities into a single string
			ragContext := ""
			for _, similarity := range similarities {
				ragContext += similarity.Prompt + " "
			}

			messages = append(messages, api.Message{Role: "system", Content: "CONTEXT:\n" + ragContext})

		} // end of similarites search

		//messages = append(messages, api.Message{Role: "user", Content: userQuestion})

		// Prompt construction
		if toolsContext != "" {
			// ✋ The result of the tools invocation is added to the user question

			if *asSystemMessage {
				messages = append(messages, api.Message{Role: "system", Content: toolsContext})
				messages = append(messages, api.Message{Role: "user", Content: userQuestion})

			} else if *asUserMessage {
				if *beforeQuestion {
					messages = append(messages, api.Message{Role: "user", Content: toolsContext})
					messages = append(messages, api.Message{Role: "user", Content: userQuestion})
				} else if *afterQuestion {
					messages = append(messages, api.Message{Role: "user", Content: userQuestion})
					messages = append(messages, api.Message{Role: "user", Content: toolsContext})

				} else {
					messages = append(messages, api.Message{Role: "user", Content: userQuestion})
					messages = append(messages, api.Message{Role: "user", Content: toolsContext})
				}
			} else {
				messages = append(messages, api.Message{Role: "user", Content: userQuestion})
				messages = append(messages, api.Message{Role: "user", Content: toolsContext})
			}

		} else {
			messages = append(messages, api.Message{Role: "user", Content: userQuestion})
		}

		if config.AddToMessages != "" { // Add the content of the file to the messages: bob --add-to-messages path/to/file
			// Add the content of the file to the messages
			addToMessages, errAdd := os.ReadFile(config.AddToMessages)
			if errAdd != nil {
				log.Fatalf("😡 Error reading add-to-messages file: %v", errAdd)
			}
			messages = append(messages, api.Message{Role: "user", Content: string(addToMessages)})
		}

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
		fmt.Println("😡 Completion error:", errCompletion)

	}

	// generate a markdown file from the value of answer
	errOutput := os.WriteFile(config.OutputPath, []byte(answer), 0644)
	if errOutput != nil {
		fmt.Println("😡 Error writing output file:", errOutput)
	}
	fmt.Println()
}
