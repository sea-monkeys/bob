# Who is this Bob?

This sample is using [qwen2.5:0.5b](https://ollama.com/library/qwen2.5:0.5b) to explain who is Bob using the `README.md` file of the repository (👋 `qwen2.5:0.5b` is very small, so check carefully the answers). 
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull qwen2.5:0.5b
```

The settings for **Bob** are stored in `.bob` directory.

Run the following command in the current directory to generate a report using the prompt in `prompt.md` and the settings in `.bob`:

```bash
bob --context ../../README.md
```

Then, check the `report.md` file for the generated content.

