Bob is a Go application that serves as a bridge between user's ideas and AI language models. It is designed to transform prompts into content. To install Bob, you can use the following command:

```bash
go install github.com/sea-monkeys/bob@latest
```

To use Bob, you can follow these steps:

1. Install the application from the GitHub repository:
```bash
go install github.com/sea-monkeys/bob@latest
```

2. Create a settings file in the current directory. The file should contain the following information:
```bash
settings.json
```

3. Create a file named `prompt.md` in the current directory. This file should contain your prompt.

4. Run the application with the following command:
```bash
bob --prompt prompt.md --settings .bob --output report.md
```

5. Bob will read your prompt, send it to Ollama with the specified configuration, and save the LLM's response to the output file.

For example, if your prompt is "Write a story about a cat", the output file will be named "story_cat.txt".