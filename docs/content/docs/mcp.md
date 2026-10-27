---
title: "MCP Server (AI Integration)"
weight: 3
---

`pd` comes with a built-in [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server, allowing AI agents like Claude Desktop, Cursor, and custom LLM workflows to interact with Pixeldrain.

## Starting the MCP Server

Start the stdio-based MCP server with:

```shell
pd mcp
```

---

## Configuration

### Claude Desktop

Add the following to your `claude_desktop_config.json` (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "pixeldrain": {
      "command": "pd",
      "args": ["mcp"],
      "env": {
        "PIXELDRAIN_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

> [!NOTE]
> `PIXELDRAIN_API_KEY` is optional. If provided, uploads will be associated with your Pixeldrain account.

### Cursor / Other MCP Clients

Configure the command as `pd` and arguments as `["mcp"]`.

---

## Available MCP Tools

The `pd mcp` server exposes the following tools:

| Tool | Description |
| :--- | :--- |
| **`upload_files`** | Uploads one or more base64-encoded files (images, audio, zip archives) and returns a file or album URL. |
| **`upload_text`** | Uploads text or markdown content as a file and returns the shareable link. |
| **`download_file`** | Downloads a file from Pixeldrain to the specified local path. |

With these tools, AI assistants can easily create sharable links for generated artifacts, charts, or images, as well as fetch remote assets during interactive sessions.
