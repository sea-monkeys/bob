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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/asellus"
	"github.com/sea-monkeys/bob/config"
	"github.com/sea-monkeys/bob/rag"
	"github.com/sea-monkeys/daphnia"
)

// TODO: check if the model is loaded / exists
// TODO: add a waiting message
// TODO: add an option for the conversational memory
// TODO: generate the report and its content at the same time (streaming)
// TODO: add several files to the messages?

/* TODO: about --create and --rag, it would be better to use command instead of flags.
Try something like this:

	// Parse command-line arguments
	flag.Parse()

	// Check command and execute
	switch flag.Arg(0) {
	case "create":
		fmt.Printf("Hello, %s!\n", *namePtr)
	case "rag":
		fmt.Println("generate vector store")
	}
*/

var (
	FALSE = false
	TRUE  = true
)

//go:embed version.txt
var versionTxt []byte

// Sample RAG files

//go:embed templates/sample.rag.env.txt
var sampleRagEnv []byte

//go:embed templates/sample.rag.instructions.txt
var sampleRagInstructions []byte

//go:embed templates/sample.rag.parameters.txt
var sampleRagParameters []byte

//go:embed templates/sample.rag.settings.txt
var sampleRagSettings []byte

//go:embed templates/sample.rag.content.txt
var sampleRagContent []byte

//go:embed templates/sample.rag.prompt.txt
var sampleRagPrompt []byte

//go:embed templates/sample.rag.readme.txt
var sampleRagReadme []byte

// Sample Schema files

//go:embed templates/sample.schema.context.txt
var sampleSchemaContext []byte

//go:embed templates/sample.schema.env.txt
var sampleSchemaEnv []byte

//go:embed templates/sample.schema.instructions.txt
var sampleSchemaInstructions []byte

//go:embed templates/sample.schema.prompt.txt
var sampleSchemaPrompt []byte

//go:embed templates/sample.schema.schema.txt
var sampleSchemaSchema []byte

//go:embed templates/sample.schema.settings.txt
var sampleSchemaSettings []byte

//go:embed templates/sample.schema.readme.txt
var sampleSchemaReadme []byte

// Sample Chat files

//go:embed templates/sample.chat.env.txt
var sampleChatEnv []byte

//go:embed templates/sample.chat.instructions.txt
var sampleChatInstructions []byte

//go:embed templates/sample.chat.prompt.txt
var sampleChatPrompt []byte

//go:embed templates/sample.chat.settings.txt
var sampleChatSettings []byte

//go:embed templates/sample.chat.readme.txt
var sampleChatReadme []byte

// Sample Tools files

//go:embed templates/sample.tools.env.txt
var sampleToolsEnv []byte

//go:embed templates/sample.tools.instructions.txt
var sampleToolsInstructions []byte

//go:embed templates/sample.tools.invocation.txt
var sampleToolsInvocation []byte

//go:embed templates/sample.tools.prompt.txt
var sampleToolsPrompt []byte

//go:embed templates/sample.tools.say_hello.txt
var sampleToolsSayHello []byte

//go:embed templates/sample.tools.settings.txt
var sampleToolsSettings []byte

//go:embed templates/sample.tools.tools.txt
var sampleToolsTools []byte

//go:embed templates/sample.tools.readme.txt
var sampleToolsReadme []byte

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

	// Create project structure
	if config.CreateProjectPathName != "" {

		// title is the last part of the path config.ProjectPathName
		title := filepath.Base(config.CreateProjectPathName)
		// The first letter must be uppercase
		title = strings.ToUpper(title[:1]) + title[1:]

		var files map[string]string
		var dirs []string

		switch kind := config.KindOfProject; kind {
		case "chat": // bob --create samples/coucou --kind chat

			dirs = []string{
				config.CreateProjectPathName,
				config.CreateProjectPathName + "/.bob",
			}

			// Define file contents
			files = map[string]string{
				filepath.Join(config.CreateProjectPathName, ".bob", ".env"):            string(sampleChatEnv),
				filepath.Join(config.CreateProjectPathName, ".bob", "instructions.md"): string(sampleChatInstructions),
				filepath.Join(config.CreateProjectPathName, ".bob", "settings.json"):   string(sampleChatSettings),
				filepath.Join(config.CreateProjectPathName, "prompt.md"):               string(sampleChatPrompt),
				filepath.Join(config.CreateProjectPathName, "README.md"):               "# " + title + "\n" + string(sampleChatReadme),
			}

		case "tools": // bob --create samples/coucou --kind tools

			dirs = []string{
				config.CreateProjectPathName,
				config.CreateProjectPathName + "/.bob",
			}

			// Define file contents
			files = map[string]string{
				filepath.Join(config.CreateProjectPathName, ".bob", ".env"):            string(sampleToolsEnv),
				filepath.Join(config.CreateProjectPathName, ".bob", "instructions.md"): string(sampleToolsInstructions),
				filepath.Join(config.CreateProjectPathName, ".bob", "settings.json"):   string(sampleToolsSettings),
				filepath.Join(config.CreateProjectPathName, ".bob", "tools.json"):      string(sampleToolsTools),
				filepath.Join(config.CreateProjectPathName, ".bob", "say_hello.sh"):    string(sampleToolsSayHello),

				filepath.Join(config.CreateProjectPathName, "tools.invocation.md"): string(sampleToolsInvocation),
				filepath.Join(config.CreateProjectPathName, "prompt.md"):           string(sampleToolsPrompt),
				filepath.Join(config.CreateProjectPathName, "README.md"):           "# " + title + "\n" + string(sampleToolsReadme),
			}

		case "rag": // bob --create samples/coucou --kind rag

			dirs = []string{
				filepath.Join(config.CreateProjectPathName, ".bob"),
				filepath.Join(config.CreateProjectPathName, "content"),
			}

			// Define files and their contents
			files = map[string]string{
				filepath.Join(config.CreateProjectPathName, ".bob", ".env"):            string(sampleRagEnv),
				filepath.Join(config.CreateProjectPathName, ".bob", "instructions.md"): string(sampleRagInstructions),
				filepath.Join(config.CreateProjectPathName, ".bob", "rag.json"):        string(sampleRagParameters),
				filepath.Join(config.CreateProjectPathName, ".bob", "settings.json"):   string(sampleRagSettings),
				filepath.Join(config.CreateProjectPathName, "content", "content.txt"):  string(sampleRagContent),
				filepath.Join(config.CreateProjectPathName, "prompt.md"):               string(sampleRagPrompt),
				filepath.Join(config.CreateProjectPathName, "README.md"):               "# " + title + "\n" + string(sampleRagReadme),
			}

		case "schema": // bob --create samples/coucou --kind schema

			dirs = []string{
				config.CreateProjectPathName,
				config.CreateProjectPathName + "/.bob",
			}

			// Define file contents
			files = map[string]string{
				filepath.Join(config.CreateProjectPathName, ".bob", ".env"):            string(sampleSchemaEnv),
				filepath.Join(config.CreateProjectPathName, ".bob", "instructions.md"): string(sampleSchemaInstructions),
				filepath.Join(config.CreateProjectPathName, ".bob", "settings.json"):   string(sampleSchemaSettings),
				filepath.Join(config.CreateProjectPathName, "context.md"):              string(sampleSchemaContext),
				filepath.Join(config.CreateProjectPathName, "prompt.md"):               string(sampleSchemaPrompt),
				filepath.Join(config.CreateProjectPathName, "README.md"):               "# " + title + "\n" + string(sampleSchemaReadme),
				filepath.Join(config.CreateProjectPathName, "schema.json"):             string(sampleSchemaSchema),
			}

		default:
			fmt.Println("🤖🤔 Kind of project not found")
		}

		// Create directories
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("😡 Error creating directory %s: %v\n", dir, err)
				return
			}
		}

		// Create and write to files
		for path, content := range files {
			if filepath.Base(path) == "say_hello.sh" {
				// Make the file executable
				err := os.WriteFile(path, []byte(content), 0755)
				if err != nil {
					fmt.Printf("😡 Error writing to file %s: %v\n", path, err)
				}
				continue
			} else {
				err := os.WriteFile(path, []byte(content), 0644)
				if err != nil {
					fmt.Printf("😡 Error writing to file %s: %v\n", path, err)
					return
				}
			}
		}

		fmt.Println("🎉 BoB project structure created successfully.")

		os.Exit(0)
	}
	// END of Project Creation

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
	/*
		if err := tools.ValidatePaths(config); err != nil {
			fmt.Printf("😡 Error: %v\n", err)
			os.Exit(1)
		}
	*/

	// Main logic
	ctx := context.Background()

	fmt.Println("🎃 config.SettingsPath", config.SettingsPath)

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
		ragConfig, errRagConf := rag.LoadRagConfig(config.SettingsPath + "/rag.json")
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
				toolsContext += string(output)
				//messages = append(messages, api.Message{Role: "system", Content: string(output)})

			}

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

			//messages = append(messages, api.Message{Role: "system", Content: toolsContext})
			//messages = append(messages, api.Message{Role: "user", Content: toolsContext})
			//userQuestion = promptContext + "\n\n" + userQuestion
		} else {

			messages = append(messages, api.Message{Role: "user", Content: userQuestion})

		}

		if config.AddToMessages != "" {
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
		log.Fatalf("😡 Completion error: %v", errCompletion)
	}

	// generate a markdown file from the value of answer
	errOutput := os.WriteFile(config.OutputPath, []byte(answer), 0644)
	if errOutput != nil {
		log.Fatalf("😡 Error writing output file: %v", errOutput)
	}
	fmt.Println()
}
