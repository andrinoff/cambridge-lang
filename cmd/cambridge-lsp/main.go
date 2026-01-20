package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/textproto"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/lexer"
	"github.com/andrinoff/cambridge-lang/pkg/parser"
	"github.com/andrinoff/cambridge-lang/pkg/token"
)

// LSP Types
const (
	TokenKeyword  = 0
	TokenString   = 1
	TokenNumber   = 2
	TokenOperator = 3
	TokenVariable = 4
	TokenComment  = 5
)

var tokenTypes = []string{
	"keyword", "string", "number", "operator", "variable", "comment",
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	documents := make(map[string]string) // Cache document content

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

		// --- INITIALIZE ---
		if method == "initialize" {
			sendResponse(request["id"], map[string]interface{}{
				"capabilities": map[string]interface{}{
					"textDocumentSync": 1, // Full sync
					"completionProvider": map[string]interface{}{
						"triggerCharacters": []string{"."},
					},
					"semanticTokensProvider": map[string]interface{}{
						"legend": map[string]interface{}{
							"tokenTypes":     tokenTypes,
							"tokenModifiers": []string{},
						},
						"range": true,
						"full":  true,
					},
				},
			})
		}

		// --- DOCUMENT SYNC ---
		if method == "textDocument/didOpen" {
			params := request["params"].(map[string]interface{})
			doc := params["textDocument"].(map[string]interface{})
			uri := doc["uri"].(string)
			text := doc["text"].(string)
			documents[uri] = text
			publishDiagnostics(uri, text)
		} else if method == "textDocument/didChange" {
			params := request["params"].(map[string]interface{})
			doc := params["textDocument"].(map[string]interface{})
			uri := doc["uri"].(string)
			changes := params["contentChanges"].([]interface{})
			if len(changes) > 0 {
				lastChange := changes[len(changes)-1].(map[string]interface{})
				text := lastChange["text"].(string)
				documents[uri] = text
				publishDiagnostics(uri, text)
			}
		}

		// --- COMPLETION ---
		if method == "textDocument/completion" {
			items := []map[string]interface{}{}

			// Add Keywords
			for k := range token.Keywords {
				items = append(items, map[string]interface{}{
					"label":  k,
					"kind":   14, // Keyword
					"detail": "keyword",
				})
			}

			// Sort for consistency
			sort.Slice(items, func(i, j int) bool {
				return items[i]["label"].(string) < items[j]["label"].(string)
			})

			sendResponse(request["id"], items)
		}

		// --- SEMANTIC TOKENS (HIGHLIGHTING) ---
		if method == "textDocument/semanticTokens/full" {
			params := request["params"].(map[string]interface{})
			docParams := params["textDocument"].(map[string]interface{})
			uri := docParams["uri"].(string)

			if text, ok := documents[uri]; ok {
				data := computeSemanticTokens(text)
				sendResponse(request["id"], map[string]interface{}{
					"data": data,
				})
			} else {
				sendResponse(request["id"], nil)
			}
		}
	}
}

func computeSemanticTokens(text string) []int {
	l := lexer.New(text)
	var data []int

	lastLine := 0
	lastStart := 0

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokenType := -1

		// Map Token Type to LSP Token Type
		if _, isKeyword := token.Keywords[tok.Literal]; isKeyword {
			tokenType = TokenKeyword
		} else {
			switch tok.Type {
			case token.INTEGER_LIT, token.REAL_LIT:
				tokenType = TokenNumber
			case token.STRING_LIT, token.CHAR_LIT:
				tokenType = TokenString
			case token.IDENT:
				tokenType = TokenVariable
			case token.ASSIGN, token.PLUS, token.MINUS, token.ASTERISK, token.SLASH,
				token.EQ, token.NOT_EQ, token.LT, token.GT, token.LT_EQ, token.GT_EQ:
				tokenType = TokenOperator
			}
		}

		if tokenType == -1 {
			continue
		}

		// Calculate LSP Delta Encoding
		// LSP uses 0-based lines and columns. Lexer provides 1-based.
		line := tok.Line - 1
		col := tok.Column - 1

		deltaLine := line - lastLine
		deltaStart := col
		if deltaLine == 0 {
			deltaStart = col - lastStart
		}

		length := len(tok.Literal)

		data = append(data, deltaLine, deltaStart, length, tokenType, 0)

		lastLine = line
		lastStart = col
	}

	return data
}

func publishDiagnostics(uri, text string) {
	l := lexer.New(text)
	p := parser.New(l)
	p.ParseProgram()
	diagnostics := []map[string]interface{}{}

	for _, errStr := range p.Errors() {
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

func sendResponse(id interface{}, result interface{}) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
	msg, _ := json.Marshal(resp)
	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(msg), msg)
}
