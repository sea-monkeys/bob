# Use Rag

This sample is using [qwen2.5:1.5b](https://ollama.com/library/qwen2.5:1.5b) for the chat completion and [snowflake-arctic-embed:33m](https://ollama.com/library/snowflake-arctic-embed:33m) for the embeddings generation.
Don't forget to pull the latest version of the models before running the script:

```bash
ollama pull qwen2.5:1.5b
ollama pull snowflake-arctic-embed:33m
```

The settings for **Bob** are stored in `.bob` directory.

You have to create a `rag.json` file in the `.bob` directory with the following content:

```json
{
    "chunkSize": 1024,
    "chunkOverlap": 256,
    "similarityThreshold": 0.3,
    "maxSimilarity": 10
}
```


## Creation (or update) of the Vector store

Run the following command in the current directory to create or update the Vector store using the content in the `content` directory:

```bash
bob --rag ./content
```
> Only text files are considered for the Vector store (md, txt, etc.)

## Completion

Run the following command in the current directory to generate a report using the prompt in `prompt.md`:

```bash
bob 

#or
bob --user "Summarize information on Gnomish Species"
```
