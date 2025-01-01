Generate a mermaid sequenceDiagram using the content below (you can use emojis to make it more fun):

## General Architecture

The Model Context Protocol uses a client-server architecture where a host application can connect to multiple servers:

1. **MCP Hosts**: These are generative AI applications using LLMs that want to access external resources via MCP. An example of a host application is Claude Desktop.
2. **MCP Clients**: Protocol clients that maintain 1:1 connections with servers (and the client is used by MCP host applications)
3. **MCP Servers**: Programs that expose specific functionalities via the MCP protocol using local or remote data sources.

The MCP protocol offers two main transport models: **STDIO** (Standard Input/Output) and **SSE** (Server-Sent Events). Both use JSON-RPC 2.0 as the message format for data transmission.

The first, **STDIO**, communicates via standard input/output streams. It is ideal for local integrations. The second, **SSE**, uses HTTP requests for communication, with SSE for server-to-client communications and POST requests for client-to-server communication. It is more suitable for remote integrations.
