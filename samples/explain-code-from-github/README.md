# Explain source code on GitHub

This sample is using [qwen2.5-coder:14b](https://ollama.com/library/qwen2.5-coder:14b) and [qwen2.5:0.5b](https://ollama.com/library/qwen2.5:0.5b).
Don't forget to pull the latest version of these modele before running the script:

```bash
ollama pull qwen2.5-coder:14b
ollama pull qwen2.5:0.5b
```

The settings for **Bob** are stored in `.bob` directory.

Run the following command in the current directory to generate a report using the tools invocation in `tools.invocation.md` and the prompt in `prompt.md` and the settings in `.bob`:

```bash
# with the default settings
bob --tools --as-user --after-question
# or
# bob --prompt prompt.md --settings .bob --output report.md --tools
```

Then, check the `report.md` file for the generated content.

