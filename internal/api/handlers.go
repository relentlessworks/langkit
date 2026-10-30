package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/langkit/internal/model"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	NoAuth bool
}

// NewHandler creates a new API handler.
func NewHandler(noAuth bool) *Handler {
	return &Handler{NoAuth: noAuth}
}

// Routes returns the HTTP mux with all routes registered.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/detect", h.detect)
	mux.HandleFunc("/batch", h.batch)
	mux.HandleFunc("/languages", h.languages)
	mux.HandleFunc("/scripts", h.scripts)
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.agentMD)
	mux.HandleFunc("/mcp", h.mcp)
	mux.HandleFunc("/", h.index)
	return mux
}

// authMiddleware checks for bearer token unless no-auth is enabled.
func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.NoAuth {
			next(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, r, http.StatusUnauthorized, "missing or invalid authorization header", "send 'Authorization: Bearer <token>' header, or run with -no-auth for dev mode")
			return
		}
		next(w, r)
	}
}

// detect handles POST /detect — detect the language of input text.
func (h *Handler) detect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use POST to detect language")
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid JSON body", "send a JSON object with a 'text' field, e.g. {\"text\": \"hello world\"}")
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		writeError(w, r, http.StatusBadRequest, "text is required", "provide non-empty text in the 'text' field")
		return
	}

	result := model.Detect(req.Text)

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, result)
		return
	}

	// Plain text: one grepable line
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("language=%s code=%s name=%s confidence=%.2f%% script=%s",
		result.Language, result.Code, result.Name, result.Confidence, result.Script))
	for _, alt := range result.Alternatives {
		sb.WriteString(fmt.Sprintf(" alt=%s(%.2f%%)", alt.Code, alt.Confidence))
	}
	writeText(w, http.StatusOK, sb.String())
}

// batch handles POST /batch — detect languages for multiple texts.
func (h *Handler) batch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use POST to batch detect languages")
		return
	}

	var req struct {
		Texts []string `json:"texts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid JSON body", "send a JSON object with a 'texts' array, e.g. {\"texts\": [\"hello\", \"bonjour\"]}")
		return
	}

	if len(req.Texts) == 0 {
		writeError(w, r, http.StatusBadRequest, "texts array is required", "provide a non-empty 'texts' array")
		return
	}

	results := model.DetectBatch(req.Texts)

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"results": results,
			"count":   len(results),
		})
		return
	}

	// Plain text: one line per result
	var sb strings.Builder
	for _, res := range results {
		sb.WriteString(fmt.Sprintf("[%d] language=%s code=%s confidence=%.2f%% script=%s\n",
			res.Index, res.Result.Language, res.Result.Code, res.Result.Confidence, res.Result.Script))
	}
	writeText(w, http.StatusOK, strings.TrimRight(sb.String(), "\n"))
}

// languages handles GET /languages — list all supported languages.
func (h *Handler) languages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use GET to list supported languages")
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"languages": model.SupportedLanguages,
			"count":     len(model.SupportedLanguages),
		})
		return
	}

	// Plain text: one line per language
	var sb strings.Builder
	for _, lang := range model.SupportedLanguages {
		sb.WriteString(fmt.Sprintf("code=%s name=%s native=%s scripts=%s\n",
			lang.Code, lang.Name, lang.NativeName, strings.Join(lang.Scripts, ",")))
	}
	writeText(w, http.StatusOK, strings.TrimRight(sb.String(), "\n"))
}

// scripts handles GET /scripts — list all detectable scripts.
func (h *Handler) scripts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "use GET to list detectable scripts")
		return
	}

	scripts := model.SupportedScripts()

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"scripts": scripts,
			"count":   len(scripts),
		})
		return
	}

	// Plain text: one line per script
	var sb strings.Builder
	for _, s := range scripts {
		sb.WriteString(s + "\n")
	}
	writeText(w, http.StatusOK, strings.TrimRight(sb.String(), "\n"))
}

// index handles GET / — service info.
func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, r, http.StatusNotFound, "endpoint not found", "GET /help for the full API manual")
		return
	}

	info := fmt.Sprintf("langkit — agentic-first language detection service\n\nEndpoints:\n  POST /detect     — detect language of text\n  POST /batch      — batch detect languages\n  GET  /languages  — list supported languages\n  GET  /scripts    — list detectable scripts\n  GET  /help       — full API manual\n  POST /mcp        — MCP JSON-RPC 2.0 endpoint\n\nSupported: %d languages, %d scripts\n\nGET /help for full documentation.",
		len(model.SupportedLanguages), len(model.SupportedScripts()))
	writeText(w, http.StatusOK, info)
}

// help handles GET /help — self-documenting API manual.
func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, helpText)
}

// agentMD handles GET /.well-known/agent.md — agent-readable manual.
func (h *Handler) agentMD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, helpText)
}
