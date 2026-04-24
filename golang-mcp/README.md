# trillium-mcp

MCP server (and companion CLI) for [TriliumNext Notes](https://github.com/TriliumNext/Notes), written in Go. Allows any MCP-compatible LLM client to read, search, create, and update notes via the Trillium ETAPI.

## Tools exposed

| Tool | Description |
|---|---|
| `search_notes` | Search notes by keyword, returns titles and content snippets |
| `get_note_content` | Fetch the full content of a note by ID (returned as Markdown) |
| `update_note` | Overwrite a note's content (accepts Markdown, converts to HTML) |
| `create_note` | Create a new note under a parent note (accepts Markdown) |

## Configuration

Create a `.env` file in the working directory:

```env
TRILLIUM_ETAPI_ADDRESS=http://your-trillium-host/etapi
TRILLIUM_ETAPI_APIKEY=your-etapi-token
GOLANG_TRILLIUM_MCP_VERSION=0.0.1
```

The ETAPI token can be generated in Trillium under **Options → ETAPI**.

## Running the MCP server

```sh
./bin/trillium-mcp
```

The server communicates over stdio and is compatible with any MCP client (e.g. Claude Desktop, Cursor).

## CLI

A development CLI is included:

```sh
# Search notes
./bin/trillium-cli search <keyword>

# Get note content
./bin/trillium-cli content --id <note-id>

# Update a note (from file or inline)
./bin/trillium-cli update --id <note-id> --path note.md
./bin/trillium-cli update --id <note-id> --content "# Hello"

# Create a note
./bin/trillium-cli add --parent <parent-id> --title "My Note" --path note.md
./bin/trillium-cli add --parent <parent-id> --title "My Note" --content "# Hello"
```

## Building

```sh
make build
```

Requires Go 1.21+ and `ogen` for code generation (`make generate`).
