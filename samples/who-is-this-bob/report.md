Bob is a Go application that serves as a bridge between your ideas and AI language models served by Ollama. It transforms your prompts into content.

To install Bob, you can use the following command:

```bash
go install github.com/sea-monkeys/bob@latest
```

To use Bob, you can follow these steps:

1. Open a terminal and navigate to the directory where you have installed Bob.
2. Run the following command to set up the environment:

```bash
bob --prompt prompt.md --settings .bob --output report.md
```

3. Bob will prompt you for a prompt and save the response to the output file.

4. Bob supports the following flags:

   - `--prompt`: Path to your prompt file.
   - `--settings`: Path to your settings file.
   - `--output`: Path to the output file.
   - `--tools`: Activate the support for tools invocation.

5. To use Bob, you can run the following command:

```bash
bob --prompt prompt.md --settings .bob --output report.md
```

This will prompt you for a prompt and save the response to the output file.

### Sample Usage

Here's a sample usage of Bob:

```bash
bob --prompt "You are a pizza expert" --settings "your-llm-settings.json" --output "pizza-expert-report.md"
```

This will prompt you for a prompt and save the response to the output file.