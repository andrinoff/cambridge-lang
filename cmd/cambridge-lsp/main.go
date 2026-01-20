package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/lexer"
	"github.com/andrinoff/cambridge-lang/pkg/parser"
)

// Minimal LSP implementation
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read Header
		header, err := textproto.NewReader(reader).ReadMIMEHeader()
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}
		length, _ := strconv.Atoi(header.Get("Content-Length"))

		// Read Body
		body := make([]byte, length)
		if _, err := io.ReadFull(reader, body); err != nil {
			continue
		}

		// Handle Request
		var request map[string]interface{}
		json.Unmarshal(body, &request)

		method, _ := request["method"].(string)

		// Initialize response
		if method == "initialize" {
			sendResponse(request["id"], map[string]interface{}{
				"capabilities": map[string]interface{}{
					"textDocumentSync": 1, // Full sync
				},
			})
		}

		// Check for errors on change or open
		if method == "textDocument/didOpen" || method == "textDocument/didChange" {
			params := request["params"].(map[string]interface{})
			var text string
			var uri string

			if method == "textDocument/didOpen" {
				doc := params["textDocument"].(map[string]interface{})
				text = doc["text"].(string)
				uri = doc["uri"].(string)
			} else {
				doc := params["textDocument"].(map[string]interface{})
				uri = doc["uri"].(string)
				changes := params["contentChanges"].([]interface{})
				if len(changes) > 0 {
					lastChange := changes[len(changes)-1].(map[string]interface{})
					text = lastChange["text"].(string)
				}
			}

			// Parse and find errors
			l := lexer.New(text)
			p := parser.New(l)
			p.ParseProgram()
			diagnostics := []map[string]interface{}{}

			for _, errStr := range p.Errors() {
				// Parse error string "line X, column Y: msg"
				// Note: This relies on your parser's specific error format
				parts := strings.SplitN(errStr, ": ", 2)
				if len(parts) < 2 {
					continue
				}

				locParts := strings.Split(parts[0], ",")
				lineStr := strings.TrimPrefix(strings.TrimSpace(locParts[0]), "line ")
				colStr := strings.TrimPrefix(strings.TrimSpace(locParts[1]), "column ")

				line, _ := strconv.Atoi(lineStr)
				col, _ := strconv.Atoi(colStr)

				diagnostics = append(diagnostics, map[string]interface{}{
					"range": map[string]interface{}{
						"start": map[string]int{"line": line - 1, "character": col - 1},
						"end":   map[string]int{"line": line - 1, "character": col + 10},
					},
					"severity": 1, // Error
					"message":  parts[1],
				})
			}

			// Publish Diagnostics
			notification := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "textDocument/publishDiagnostics",
				"params": map[string]interface{}{
					"uri":         uri,
					"diagnostics": diagnostics,
				},
			}
			msg, _ := json.Marshal(notification)
			fmt.Printf("Content-Length: %d\r\n\r\n%s", len(msg), msg)
		}
	}
}

func sendResponse(id interface{}, result interface{}) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
	msg, _ := json.Marshal(resp)
	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(msg), msg)
}
