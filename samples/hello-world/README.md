# Hello World


This sample is using [deepseek-coder:1.3b](https://ollama.com/library/deepseek-coder:1.3b) to generate hello world programs in various language. 
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull deepseek-coder:1.3b
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

