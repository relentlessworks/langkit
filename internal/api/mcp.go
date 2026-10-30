package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/relentlessworks/langkit/internal/model"
)

// mcp handles POST /mcp — Model Context Protocol JSON-RPC 2.0 endpoint.
func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use POST for MCP JSON-RPC")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, mcpError(nil, -32700, "Parse error"))
		return
	}

	method, _ := req["method"].(string)
	id := req["id"]

	switch method {
	case "initialize":
		writeJSON(w, http.StatusOK, mcpResult(id, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "langkit",
				"version": "0.1.0",
			},
		}))

	case "tools/list":
		writeJSON(w, http.StatusOK, mcpResult(id, map[string]interface{}{
			"tools": []map[string]interface{}{
				{
					"name":        "detect_language",
					"description": "Detect the language of text. Returns language code, name, confidence, and script.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"text": map[string]interface{}{
								"type":        "string",
								"description": "The text to analyze",
							},
						},
						"required": []string{"text"},
					},
				},
				{
					"name":        "batch_detect",
					"description": "Detect languages for multiple texts at once.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"texts": map[string]interface{}{
								"type":        "array",
								"items":       map[string]interface{}{"type": "string"},
								"description": "Array of texts to analyze",
							},
						},
						"required": []string{"texts"},
					},
				},
			},
		}))

	case "tools/call":
		params, _ := req["params"].(map[string]interface{})
		toolName, _ := params["name"].(string)
		args, _ := params["arguments"].(map[string]interface{})

		switch toolName {
		case "detect_language":
			text, _ := args["text"].(string)
			if text == "" {
				writeJSON(w, http.StatusOK, mcpError(id, -32602, "Missing required parameter: text"))
				return
			}
			result := model.Detect(text)
			writeJSON(w, http.StatusOK, mcpResult(id, map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf("language=%s code=%s name=%s confidence=%.2f%% script=%s",
							result.Language, result.Code, result.Name, result.Confidence, result.Script),
					},
				},
			}))

		case "batch_detect":
			textsRaw, _ := args["texts"].([]interface{})
			if len(textsRaw) == 0 {
				writeJSON(w, http.StatusOK, mcpError(id, -32602, "Missing required parameter: texts"))
				return
			}
			texts := make([]string, len(textsRaw))
			for i, t := range textsRaw {
				texts[i], _ = t.(string)
			}
			results := model.DetectBatch(texts)
			var sb string
			for _, res := range results {
				sb += fmt.Sprintf("[%d] language=%s code=%s confidence=%.2f%% script=%s\n",
					res.Index, res.Result.Language, res.Result.Code, res.Result.Confidence, res.Result.Script)
			}
			writeJSON(w, http.StatusOK, mcpResult(id, map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": sb,
					},
				},
			}))

		default:
			writeJSON(w, http.StatusOK, mcpError(id, -32601, fmt.Sprintf("Unknown tool: %s", toolName)))
		}

	default:
		writeJSON(w, http.StatusOK, mcpError(id, -32601, fmt.Sprintf("Unknown method: %s", method)))
	}
}

func mcpResult(id interface{}, result interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}

func mcpError(id interface{}, code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}
