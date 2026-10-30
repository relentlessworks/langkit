package api

const helpText = `# LangKit — Agentic-First Language Detection Service

LangKit detects the language of text using Unicode script analysis, common word
matching, and statistical scoring. 70+ languages across 20+ writing systems.

## Quick Start

    curl -X POST http://localhost:7700/detect -d '{"text":"Hello world"}'

Response (plain text):
    language=en code=en name=English confidence=85.50% script=Latin

JSON mode:
    curl -X POST http://localhost:7700/detect -H "Accept: application/json" -d '{"text":"Bonjour le monde"}'

## Endpoints

### POST /detect — Detect Language
Detect the language of a single text.

Request body (JSON):
    {"text": "The text to analyze"}

Response (plain text):
    language=en code=en name=English confidence=85.50% script=Latin alt=de(10.20%)

Response (JSON):
    {"language":"en","code":"en","name":"English","native_name":"English",
     "confidence":85.5,"script":"Latin","alternatives":[...]}

### POST /batch — Batch Detect
Detect languages for multiple texts at once.

Request body (JSON):
    {"texts": ["Hello world", "Bonjour le monde", "Hola mundo"]}

Response (plain text, one line per text):
    [0] language=en code=en confidence=85.50% script=Latin
    [1] language=fr code=fr confidence=72.30% script=Latin
    [2] language=es code=es confidence=78.10% script=Latin

### GET /languages — List Supported Languages
Returns all 70+ supported languages with codes, names, and scripts.

### GET /scripts — List Detectable Scripts
Returns all 20+ Unicode scripts langkit can identify.

### POST /mcp — MCP JSON-RPC 2.0
Model Context Protocol endpoint for chat client integrations.
Supports: tools/list, tools/call (detect_language, batch_detect).

## Output Format

- **Plain text** (default): One grepable line per record. Token-cheap.
- **JSON**: Add Accept: application/json or ?format=json to any request.

## Auth

- **Dev mode**: Run with -no-auth flag to skip auth.
- **Production**: Send Authorization: Bearer <token> header.

## Configuration

    LANGKIT_ADDR=:7700     Listen address (default :7700)
    LANGKIT_NO_AUTH=true   Disable auth (dev mode)
    LANGKIT_SECRET=...     Token signing secret (auto-generated if empty)

    Flags: -addr, -no-auth, -secret

## Detection Method

LangKit uses a multi-stage approach:
1. **Script detection** — Identifies the Unicode script (Latin, Cyrillic, Arabic, etc.)
2. **Common word matching** — Matches frequent words per language
3. **Statistical scoring** — Normalizes scores and ranks candidates
4. **Confidence** — Percentage based on relative word match frequency

## Error Handling

All 4xx errors include a "hint" field telling you what to do next:
    error: text is required | hint: provide non-empty text in the 'text' field
`
