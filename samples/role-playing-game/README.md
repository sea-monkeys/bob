# Role Playing Game


This sample is using [nemotron-mini:4b](https://ollama.com/library/nemotron-mini:4b) to generate a role-playing game scenario.
Don't forget to pull the latest version of the model before running the script:

```bash
ollama pull nemotron-mini:4b
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

