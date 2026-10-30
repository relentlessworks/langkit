package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *httptest.Server {
	h := NewHandler(true) // no-auth mode
	return httptest.NewServer(h.Routes())
}

func TestDetectEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"text":"The quick brown fox jumps over the lazy dog"}`)
	req, _ := http.NewRequest("POST", server.URL+"/detect", body)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["code"] != "en" {
		t.Errorf("expected en, got %v", result["code"])
	}
}

func TestDetectEndpointPlainText(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"text":"The quick brown fox jumps over the lazy dog"}`)
	resp, err := http.Post(server.URL+"/detect", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	text := string(respBody)
	if !strings.Contains(text, "language=en") {
		t.Errorf("expected plain text with language=en, got: %s", text)
	}
}

func TestDetectEndpointJSON(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"text":"Bonjour le monde, ceci est un texte en français pour la détection"}`)
	req, _ := http.NewRequest("POST", server.URL+"/detect", body)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["code"] != "fr" {
		t.Errorf("expected fr, got %v", result["code"])
	}
}

func TestDetectEmptyText(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"text":""}`)
	resp, err := http.Post(server.URL+"/detect", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDetectInvalidJSON(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`not json`)
	resp, err := http.Post(server.URL+"/detect", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestBatchEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"texts":["Hello world this is English","Bonjour le monde ceci est français","Hola mundo este es español"]}`)
	req, _ := http.NewRequest("POST", server.URL+"/batch", body)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	results, _ := result["results"].([]interface{})
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestLanguagesEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/languages", nil)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	langs, _ := result["languages"].([]interface{})
	if len(langs) < 60 {
		t.Errorf("expected at least 60 languages, got %d", len(langs))
	}
}

func TestScriptsEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/scripts", nil)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	scripts, _ := result["scripts"].([]interface{})
	if len(scripts) < 20 {
		t.Errorf("expected at least 20 scripts, got %d", len(scripts))
	}
}

func TestHelpEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/help")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIndexEndpoint(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestNotFound(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestMCPInitialize(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`)
	resp, err := http.Post(server.URL+"/mcp", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	res, _ := result["result"].(map[string]interface{})
	if res == nil {
		t.Fatal("expected result in response")
	}
}

func TestMCPToolsList(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	resp, err := http.Post(server.URL+"/mcp", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	res, _ := result["result"].(map[string]interface{})
	tools, _ := res["tools"].([]interface{})
	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}
}

func TestMCPDetectLanguage(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	body := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"detect_language","arguments":{"text":"Hello world this is English"}}}`)
	resp, err := http.Post(server.URL+"/mcp", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	res, _ := result["result"].(map[string]interface{})
	content, _ := res["content"].([]interface{})
	if len(content) == 0 {
		t.Fatal("expected content in response")
	}
}

func TestDetectMethodNotAllowed(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/detect")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}
