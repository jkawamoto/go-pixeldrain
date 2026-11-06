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

The server writes lifecycle and tool execution logs to standard error. Tool arguments and results are not logged; the MCP client determines whether standard error is displayed or saved.

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
| **`list_recipients`** | Lists the configured public key IDs for encrypted uploads. |

With these tools, AI assistants can easily create sharable links for generated artifacts, charts, or images, as well as fetch remote assets during interactive sessions.

## Encrypted transfers

Configure the public key directories when starting the server:

```shell
pd mcp --age-recipients-dir /path/to/age-keys --ssh-recipients-dir /path/to/ssh-keys
```

Both options are optional. Each directory contains one public key per `.pub` file; other files and symbolic links are ignored. The age directory accepts age X25519 public keys, and the SSH directory accepts supported SSH public keys. Blank lines and `#` comments are ignored. For example, `alice.pub` in the age directory appears in `list_recipients` as `age:alice.pub`, with the display name `alice`. SSH IDs use the `ssh:` prefix.

Call `list_recipients` and pass the chosen IDs as the `recipients` array to `upload_text` or `upload_files`. Each file is encrypted for all selected recipients and uploaded with `.age` appended to its name. Omitting `recipients` preserves the unencrypted behavior; an empty array, invalid ID, or invalid key is an error.

To decrypt `.age` downloads, start the server with `pd mcp --identity-file /path/to/key`, using an age identity file or an unencrypted SSH private key. The identity is loaded at startup; restart the server after changing it. Without this option, encrypted files are downloaded unchanged. Downloads do not overwrite existing files; choose a different output path if one exists. Passphrase-protected SSH keys require an interactive terminal and are not supported by the MCP server. The identity path may be visible in process listings or MCP client configuration, but is not an MCP tool argument.

The file name and approximate encrypted size remain visible to Pixeldrain. Text passed through MCP and recipient IDs may also be retained by the MCP client or in conversation history. Treat recipient names as untrusted metadata and confirm ambiguous recipients before uploading.
