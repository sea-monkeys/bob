This Go program is a command-line tool that interacts with an Ollama API to perform various tasks such as chat completion, embedding generation, and retrieval-augmented generation (RAG). The program is designed to be flexible and configurable through command-line flags and configuration files. Below is a detailed breakdown of the program's functionality and components:

### Key Components

1. **Command-Line Flags**:
   - `--model`: Specifies the model to use for chat completion.
   - `--embeddings-model`: Specifies the model to use for embedding generation.
   - `--json-schema`: Indicates whether to use a JSON schema for the response format.
   - `--tools`: Indicates whether to use tools for the response.
   - `--rag`: Indicates whether to use retrieval-augmented generation.
   - `--tools-config`: Specifies the path to the tools configuration file.
   - `--rag-config`: Specifies the path to the RAG configuration file.
   - `--json-schema-path`: Specifies the path to the JSON schema file.
   - `--settings-path`: Specifies the path to the settings directory.
   - `--output-path`: Specifies the path to the output file.
   - `--ollama-raw-url`: Specifies the URL of the Ollama API.

2. **Configuration Files**:
   - `settings.json`: Contains configuration settings for the program.
   - `rag.json`: Contains configuration settings for retrieval-augmented generation.
   - `schema.json`: Contains a JSON schema for the response format.

3. **Main Functionality**:
   - **Initialization**: The program initializes various components such as the Ollama client, configuration settings, and vector store.
   - **Embedding Generation**: If the `--embeddings-model` flag is set, the program generates embeddings for the input question.
   - **Retrieval-Augmented Generation (RAG)**: If the `--rag` flag is set, the program searches for similar chunks in a vector store and uses them to augment the response.
   - **Chat Completion**: The program sends a chat completion request to the Ollama API and processes the response.
   - **Output**: The final response is written to an output file in Markdown format.

### Detailed Steps

1. **Initialization**:
   - The program starts by parsing command-line flags and loading configuration files.
   - It initializes the Ollama client with the specified URL.
   - It loads the settings from the `settings.json` file and sets up the vector store if necessary.

2. **Embedding Generation**:
   - If the `--embeddings-model` flag is set, the program generates embeddings for the input question using the specified embeddings model.
   - The embeddings are used to search for similar chunks in the vector store.

3. **Retrieval-Augmented Generation (RAG)**:
   - If the `--rag` flag is set, the program searches for similar chunks in the vector store based on the generated embeddings.
   - The similar chunks are used to augment the context for the chat completion request.

4. **Chat Completion**:
   - The program constructs a chat completion request with the appropriate messages and options.
   - If the `--json-schema` flag is set, the program uses a JSON schema for the response format.
   - The program sends the chat completion request to the Ollama API and processes the response.

5. **Output**:
   - The final response is written to an output file in Markdown format.

### Example Usage

```sh
go run main.go --model "gpt-4" --embeddings-model "text-embedding-ada-002" --rag --settings-path "./settings" --output-path "./output.md"
```

This command would run the program with the specified model, embeddings model, and RAG settings, and write the output to `output.md`.

### Conclusion

This program is a versatile tool for interacting with the Ollama API, providing functionalities such as chat completion, embedding generation, and retrieval-augmented generation. It is designed to be highly configurable and can be adapted to various use cases by modifying the configuration files and command-line flags.