package config

type Config struct {
	PromptPath          string
	ContextPath         string // for this one check if the file exists
	ToolsInvocationPath string
	JsonSchemaPath      string

	SettingsPath     string
	OutputPath       string
	RagDocumentsPath string // for RAG

	// Generate a project structure
	CreateProjectPathName string
	KindOfProject         string

	// to override the system and user questions
	System string
	User   string

	// --add-to-messages ../../main.go --as-user --after-question
	AddToMessages string
}
