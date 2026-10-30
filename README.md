# LangKit

> Agentic-first language detection service. No UI, no SDK — the agent IS the interface.

LangKit detects the language of text using Unicode script analysis, common word
matching, and statistical scoring. **70+ languages** across **20+ writing systems**.

## Quick Start

```bash
# Build
make build

# Run (dev mode, no auth)
./langkit -no-auth

# Detect
curl -X POST http://localhost:7700/detect -d '{"text":"Hello world"}'
```

Output (plain text, one grepable line):
```
language=en code=en name=English confidence=85.50% script=Latin
```

JSON mode:
```bash
curl -X POST http://localhost:7700/detect -H "Accept: application/json" -d '{"text":"Bonjour le monde"}'
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/detect` | Detect language of text |
| POST | `/batch` | Batch detect languages |
| GET | `/languages` | List supported languages |
| GET | `/scripts` | List detectable scripts |
| GET | `/help` | Full API manual |
| GET | `/.well-known/agent.md` | Agent-readable manual |
| POST | `/mcp` | MCP JSON-RPC 2.0 endpoint |

## Design Principles

1. **The agent IS the interface** — No UI, no SDK. Plain HTTP.
2. **Plain text by default** — One labeled, grepable line per record. JSON on demand.
3. **Instructive errors** — Every 4xx includes a hint telling the agent what to do next.
4. **Self-documenting** — `GET /help` returns a one-page operating manual.
5. **Single static binary** — Go, zero external dependencies, CGO_ENABLED=0.
6. **Zero config defaults** — Runs out of the box.
7. **MCP connector** — Speaks Model Context Protocol at `/mcp`.

## Detection Method

LangKit uses a multi-stage approach:

1. **Script detection** — Identifies the Unicode script (Latin, Cyrillic, Arabic, etc.)
2. **Common word matching** — Matches frequent words per language
3. **Statistical scoring** — Normalizes scores and ranks candidates
4. **Confidence** — Percentage based on relative word match frequency

## Configuration

```bash
# Environment variables
LANGKIT_ADDR=:7700       # Listen address
LANGKIT_NO_AUTH=true      # Disable auth (dev mode)
LANGKIT_SECRET=...        # Token signing secret

# Flags
./langkit -addr=:7700 -no-auth -secret=mysecret
```

## Build

```bash
make build    # Build the binary
make test     # Run tests
make vet      # Run go vet
make all      # vet + test + build
```

## License

MIT
