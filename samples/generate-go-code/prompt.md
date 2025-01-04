Generate a Golang code to generate this file structure:

```bash
samples/rust-tutorial
├── .bob
│   ├── .env
│   ├── instructions.md
│   ├── rag.json
│   └── settings.json
├── content
│   └── content.txt
├── prompt.md
└── README.md
```

- .env is a file
- instructions.md is a file
- rag.json is a file
- settings.json is a file
- content is a directory
- content.txt is a file
- prompt.md is a file
- README.md is a file

The content of .env is:
```env
OLLAMA_HOST=http://localhost:11434
LLM=qwen2.5:1.5b
EMBEDDINGS_LLM=snowflake-arctic-embed:33m
```

The content of instructions.md is:
```markdown
You are a useful AI agent.
```

The content of settings.json is:
```json
{
    "temperature": 0.8,
    "repeat_last_n": 2
}
```

The content of rag.json is:
```json
{
    "chunkSize": 1024,
    "chunkOverlap": 256,
    "similarityThreshold": 0.3,
    "maxSimilarity": 10
}
```

The content of prompt.md is:
```markdown
Put your prompt here.
```

The content of README.md is:
```markdown
# The name of the project
```