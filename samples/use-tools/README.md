# Use Tools

This sample is using [qwen2.5:1.5b](https://ollama.com/library/qwen2.5:1.5b).
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull qwen2.5:1.5b
```

The settings for **Bob** are stored in `.bob` directory.

Run the following command in the current directory to generate a report using the tools invocation in `tools.invocation.md` and the prompt in `prompt.md` and the settings in `.bob`:

```bash
# with the default settings
bob --tools
# or
# bob --tools --tools-invocation tools.invocation.md  --prompt prompt.md --settings .bob --output report.md
```
Then, check the `report.md` file for the generated content.

