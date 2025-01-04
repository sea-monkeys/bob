Bob is a powerful tool that can help you transform your prompts into content. It's a Go application that works as a bridge between your ideas and AI language models served by Ollama. Here's how you can use Bob:

1. **Requirements:**
   - Ollama installed and running locally (default: http://localhost:11434)
   - The model specified in settings (default: qwen2.5:0.5b)

2. **Directory structure needed:**
   ```
   .bob/
     ├── .env            # Contains OLLAMA_HOST and LLM settings
     ├── settings.json   # Model configuration
     └── instructions.md # System instructions for the LLM
   ```

3. **Basic usage:**
   ```bash
   bob --prompt prompt.md --settings .bob --output report.md
   ```

4. **Optional flags:**
   - `--version`: Shows version info
   - `--prompt`: Path to prompt file (default: prompt.md)
   - `--context`: Path to context file (default: context.md)
   - `--settings`: Path to settings directory (default: .bob)
   - `--output`: Path to output file (default: report.md)
   - `--tools`: Activate the support for tools invocation
   - `--schema`: Extract data from the generated content (JSON structured output)
   - `--rag`: Create or update a Vector store from a content directory

5. **Tools invocation:**
   ```bash
   bob --tools --as-user --after-question
   ```

   - `bob --tools`: Add the result of the tool(s) invocation as a `user` message.
   - `bob --tools --as-user`: Add the result of the tool(s) invocation as a `user` message, before the user's question.
   - `bob --tools --as-system`: Add the result of the tool(s) invocation as a `system` message, before the user's question.

6. **Structured JSON output:**
   ```bash
   bob --create path/to/project --kind chat|rag|tools|schema
   ```

   - `bob --create`: Create a new project structure.
   - `--kind`: Specify the project structure type (chat|rag|tools|schema).
   - `--schema`: Add the result of the tool(s) invocation as a `system` message, before the user's question.

7. **RAG support:**
   ```bash
   bob --create path/to/project --kind chat|rag|tools|schema
   ```

   - `bob --create`: Create a new project structure.
   - `--kind`: Specify the project structure type (chat|rag|tools|schema).
   - `--rag`: Add the result of the tool(s) invocation as a `system` message, before the user's question.

By following these steps, you can use Bob to transform your prompts into content, add tools invocation, use structured JSON output, and add RAG support. Let me know if you have any questions!