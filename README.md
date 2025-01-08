# BoB
> We are legion

🤖 Meet Bob: Your AI Writing Companion:

Bob is a sleek Go application that works as a bridge between your ideas and AI language models served by Ollama to transform your prompts into content.

## Installation

### Install from github

```bash
go install github.com/sea-monkeys/bob@latest
```

### Install from source

```bash
git clone https://github.com/sea-monkeys/bob.git
cd bob
go build
sudo mv bob /usr/local/bin/
```

## How to use Bob

1. Requirements:
- Ollama installed and running locally (default: http://localhost:11434)
- The model specified in settings (default: qwen2.5:0.5b)

2. Directory structure needed:
```
.bob/
  ├── .env            # Contains OLLAMA_HOST and LLM settings
  ├── settings.json   # Model configuration
  └── instructions.md # System instructions for the LLM
```

3. Basic usage:
```bash
bob --prompt prompt.md --settings .bob --output report.md

bob --prompt samples/json-to-go-struct/prompt.md --settings samples/json-to-go-struct/.bob --output ./report.md
```

4. Optional flags:
```bash
--version # Shows version info
--prompt  # Path to prompt file (default: prompt.md)
--context # Path to context file (default: context.md)
--settings # Path to settings directory (default: .bob)
--output   # Path to output file (default: report.md)
--tools    # Activate the support for tools invocation
--schema   # Extract data from the generated content (JSON structured output)
--rag      # Create or update a Vector store from a content directory
```

> You can override the content of the `prompt.md` and the `instructions.md` files with the `--system` and `--user` flags: `bob --system "you are a pizza expert" --user "what is the best pizza in the world?"`.

The application reads your prompt, sends it to Ollama with the specified configuration, and saves the LLM's response to the output file.

## Tools invocation

See the samples in the `samples` directory for examples of how to use tools invocation:

- [Use Tools](samples/use-tools)
- [Summarize web page](samples/summarize-web-page)
- [Explain code from GitHub](samples/explain-code-from-github)

If you want to use tools invocation, you need to add the `--tools` flag to the command.

If you want to add the result of the tool(s) invocation as a `user` message, after the user's question, use this command.
```bash
bob --tools --as-user --after-question
# or
bob --tools --as-user
# or (default)
bob --tools
```

If you want to add the result of the tool(s) invocation as a `user` message, before the user's question, use this command.
```bash
bob --tools --as-user --before-question
```

If you want to add the result of the tool(s) invocation as a `system` message, use this command.
```bash
bob --tools --as-system
```
> The result will be added as a `system` message before the user question to the messages list.



## Structured JSON output

See the samples in the `samples` directory for examples of how to use structured JSON output:

- [Extract data](samples/extract-data)

## RAG support

See the samples in the `samples` directory for examples of how to use RAG support:

- [Chronicles of Aethelgard](samples/chronicles-of-aethelgard)

## Generate a BoB project structure

```bash
bob --create path/to/project --kind chat|rag|tools|schema
```

## Add the content of a file at the end of the messages list

```bash
bob --add-to-messages path/to/file
```

## Messages list construction

The messages list is constructed as follows:
- system message (instructions: from instructions.md (in the settingsPath) or from the --system flag to define a string value, but not a path)
- system message (from the --context flag to define a path to context.md)
- user message (prompt/question: from prompt.md (in the settingsPath) or from the --user flag to define a string value, but not a path)
- user message (from the ----add-to-messages path/to/file)
