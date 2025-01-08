package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/bob/config"
)

func ToolsInvocation(ctx context.Context, config config.Config, ollamaClient *api.Client, ollamaRawUrl, toolsModel string) (string, error) {
	var (
		FALSE = false
		//TRUE  = true
	)

	toolsContext := ""
	// Tool invocation
	fmt.Println("🛠️🤖 using:", ollamaRawUrl, toolsModel, "for tools")

	// Read tools
	toolsConfigFile, errToolsConf := os.ReadFile(config.SettingsPath + "/tools.json")

	if errToolsConf != nil {
		fmt.Println("😡 Error reading tools.json file:", errToolsConf)
		return "", errToolsConf
	}
	var toolsList api.Tools
	errJsonToolsConf := json.Unmarshal(toolsConfigFile, &toolsList)
	if errJsonToolsConf != nil {
		fmt.Println("😡 Error unmarshalling tools.json file:", errJsonToolsConf)
		return "", errJsonToolsConf
	}

	// Load the content of the tools.invocation.md file
	toolsPrompt, errPrompt := os.ReadFile(config.ToolsInvocationPath)

	if errPrompt != nil {
		fmt.Println("😡 Error reading tools.invocation file:", errPrompt)
		return "", errPrompt
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
				fmt.Println("😡 Error executing bash tool:", err)
				//panic(err)
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
		fmt.Println("😡 Error when executing tools with Ollama", err)
		return "", err
	}

	return toolsContext, nil
}
