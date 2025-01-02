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
	"strings"

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
)

// --- MCP ---
type MCPConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// --- MCP ---

type Config struct {
	PromptPath          string
	ToolsInvocationPath string
	SettingsPath        string
	OutputPath          string
	//Version string
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

	flag.StringVar(&config.ToolsInvocationPath, "tools-invocation", "tools.invocation.md", "Path to tools invocation file")

	flag.StringVar(&config.System, "system", "", "System instructions")
	flag.StringVar(&config.User, "user", "", "User question")

	// Version flag
	version := flag.Bool("version", false, "Display version information")

	toolsInvocation := flag.Bool("tools", false, "Tools invocation")

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
	// TODO: check if the model is loaded / exists
	// TODO: add a waiting message
	// TODO: add an option for the conversational memory
	// TODO: add RAG features
	// TODO: generate the report and its content at the same time (streaming)

	url, _ := url.Parse(ollamaRawUrl)

	fmt.Println("🤖 using:", ollamaRawUrl, model)

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
	// Tools
	// ==========================================================
	promptContext := "<documents>"

	if *toolsInvocation {
		// Tool invocation
		//fmt.Println("🙂 Tool invocation not implemented yet")

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
				cmdArgs := []string{config.SettingsPath + "/" + toolCall.Function.Name+".sh"}
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
				promptContext += "<document>"+string(output)+"</document>"

				//messages = append(messages, api.Message{Role: "system", Content: string(output)})

			}
			promptContext += "</documents>"
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
	if promptContext != "" {
		userQuestion = promptContext + "\n\n" + userQuestion
	}
	messages = append(messages, api.Message{Role: "user", Content: userQuestion})

	/*
	messages = []api.Message{
		{Role: "system", Content: systemInstructions},
		{Role: "user", Content: userQuestion},
	}
	*/


	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options:  modelConfig,
		Stream:   &TRUE,
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
