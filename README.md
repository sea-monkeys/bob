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
```

4. Optional flags:
```bash
--version # Shows version info
--prompt  # Path to prompt file (default: prompt.md)
--settings # Path to settings directory (default: .bob)
--output   # Path to output file (default: report.md)
```

The application reads your prompt, sends it to Ollama with the specified configuration, and saves the LLM's response to the output file.
