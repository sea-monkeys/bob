# Generate source code from files and directory structure

This sample is using [qwen2.5-coder:0.5b](https://ollama.com/library/qwen2.5-coder:0.5b) to generate source code. 
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull qwen2.5-coder:0.5b
```


The settings for **Bob** are stored in `.bob` directory.

Run the following command in the current directory to generate a report using the prompt in `prompt.md` and the settings in `.bob`:

```bash
# with the default settings
bob
# or
# bob --prompt prompt.md --settings .bob --output report.md
```

Then, check the `report.md` file for the generated content.

