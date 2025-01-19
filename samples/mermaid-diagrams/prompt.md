Generate a mermaid sequenceDiagram using the source code below (you can use emojis to make it more fun):


```golang

func main() {

	ctx := context.Background()

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	model := os.Getenv("LLM")

	fmt.Println("🌍", ollamaUrl, "📕", model)

	client, errCli := api.ClientFromEnvironment()
	if errCli != nil {
		log.Fatal("😡:", errCli)
	}

	systemInstructions, err := os.ReadFile("instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	generationInstructions, err := os.ReadFile("steps.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Get the character
	character, errChar := GetCharacter()
	if errChar != nil {
		log.Fatal("😡:", errChar)
	}

	fmt.Println("🧙‍♂️", character.Name, "🧝‍♂️", character.Kind)

	userContent := fmt.Sprintf("Using the steps below, create a %s with this name:%s", character.Kind, character.Name)

	// Prompt construction
	messages := []api.Message{
		{Role: "system", Content: string(systemInstructions)},
		//{Role: "system", Content: string(generationInstructions)},
		{Role: "user", Content: userContent},
		{Role: "user", Content: string(generationInstructions)},
	}

	stream := true
	//noStream  := false

	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options: map[string]interface{}{
			//"temperature":   0.0,
			"temperature":    1.0,
			"repeat_last_n":  2,
			"repeat_penalty": 2.0,
			"top_k":          10,
			"top_p":          0.5,

			//"num_ctx":       4096, // https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
		},
		//Format:    "json",
		KeepAlive: &api.Duration{Duration: 1 * time.Minute},
		Stream:    &stream,
	}

	mdResult := ""
	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		mdResult += resp.Message.Content
		return nil
	}

	// Start the chat completion
	errChat := client.Chat(ctx, req, respFunc)
	if errChat != nil {
		log.Fatal("😡:", errChat)
	}

	// Character sheet
	characterSheetId := strings.ToLower(strings.ReplaceAll(character.Name, " ", "-"))

	log.Printf("Attempting to write file: ./character-sheet-%s.md", characterSheetId)

	errWriteFile := os.WriteFile("./character-sheet-"+characterSheetId+".md", []byte("# CHARACTER SHEET\n\n"+mdResult), 0644)
	if errWriteFile != nil {
		log.Fatal("😡:", errChat)
	}

	fmt.Println("\n📝", characterSheetId, "saved.")

	fmt.Println("\n🟦")
	for {
		// Loop forever
	}
}

```