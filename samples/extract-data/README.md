# Extract data


This sample is using [granite3-moe:1b](https://ollama.com/library/granite3-moe:1b).
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull granite3-moe:1b
```

The settings for **Bob** are stored in `.bob` directory.

Run the following command in the current directory to generate a report using the prompt in `prompt.md` and the settings in `.bob` and extract the data:

```bash
# with the default settings
bob --schema
# or
# bob --prompt prompt.md --settings .bob --output report.md --schema
```

Then, check the `report.md` file for the generated content.

