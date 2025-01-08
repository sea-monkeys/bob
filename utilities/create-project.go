package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sea-monkeys/bob/config"

	_ "embed"
)

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

func CreateProject(config config.Config) error {

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
			return err
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
				return err
			}
		}
	}

	fmt.Println("🎉 BoB project structure created successfully.")

	return nil
}
