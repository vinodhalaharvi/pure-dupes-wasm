// mcp-server.go - MCP Server for pure-dupes
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// MCP Protocol types
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Main MCP Server
func main() {
	log.SetOutput(os.Stderr)
	log.Println("🔍 pure-dupes MCP Server starting...")

	// Read from stdin, write to stdout
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		log.Printf("Received request: %s", req.Method)

		var resp MCPResponse
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "initialize":
			resp.Result = handleInitialize()

		case "tools/list":
			resp.Result = handleToolsList()

		case "tools/call":
			result, err := handleToolCall(req.Params)
			if err != nil {
				resp.Error = &MCPError{
					Code:    -32603,
					Message: err.Error(),
				}
			} else {
				resp.Result = result
			}

		default:
			resp.Error = &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			}
		}

		if err := encoder.Encode(resp); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}

	log.Println("MCP Server shutting down")
}

func handleInitialize() map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]string{
			"name":    "pure-dupes",
			"version": "1.0.0",
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
	}
}

func handleToolsList() map[string]interface{} {
	tools := []Tool{
		{
			Name:        "analyze_duplicates",
			Description: "Analyze a directory for duplicate files using Merkle tree-based content hashing. Finds exact and partial duplicates.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"directory": map[string]interface{}{
						"type":        "string",
						"description": "Path to directory to analyze",
					},
					"threshold": map[string]interface{}{
						"type":        "number",
						"description": "Similarity threshold (0.0-1.0) for partial matches. Default: 0.8",
						"default":     0.8,
					},
					"max_depth": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum directory depth to scan. Default: 10",
						"default":     10,
					},
				},
				"required": []string{"directory"},
			},
		},
		{
			Name:        "get_duplicate_groups",
			Description: "Get smart duplicate groups showing files that can be safely removed",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"directory": map[string]interface{}{
						"type":        "string",
						"description": "Path to directory previously analyzed",
					},
				},
				"required": []string{"directory"},
			},
		},
		{
			Name:        "check_file_hash",
			Description: "Get content hash for a specific file",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to file to hash",
					},
				},
				"required": []string{"file_path"},
			},
		},
	}

	return map[string]interface{}{
		"tools": tools,
	}
}

func handleToolCall(paramsRaw json.RawMessage) (interface{}, error) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(paramsRaw, &params); err != nil {
		return nil, fmt.Errorf("invalid params: %v", err)
	}

	log.Printf("Tool call: %s", params.Name)

	switch params.Name {
	case "analyze_duplicates":
		return analyzeDirectoryTool(params.Arguments)

	case "get_duplicate_groups":
		return getDuplicateGroupsTool(params.Arguments)

	case "check_file_hash":
		return checkFileHashTool(params.Arguments)

	default:
		return nil, fmt.Errorf("unknown tool: %s", params.Name)
	}
}

func analyzeDirectoryTool(args map[string]interface{}) (interface{}, error) {
	directory, ok := args["directory"].(string)
	if !ok {
		return nil, fmt.Errorf("directory parameter required")
	}

	threshold := 0.8
	if t, ok := args["threshold"].(float64); ok {
		threshold = t
	}

	maxDepth := 10
	if d, ok := args["max_depth"].(float64); ok {
		maxDepth = int(d)
	}

	log.Printf("Analyzing directory: %s (threshold: %.2f, depth: %d)", directory, threshold, maxDepth)

	// TODO: Implement actual directory analysis
	// For now, return mock data
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf(
					"Analyzing directory: %s\n\nConfiguration:\n- Threshold: %.0f%%\n- Max Depth: %d\n\nThis would scan the directory using Merkle trees to find duplicates.\n\nNote: Full implementation requires file system access.",
					directory, threshold*100, maxDepth,
				),
			},
		},
	}, nil
}

func getDuplicateGroupsTool(args map[string]interface{}) (interface{}, error) {
	directory, ok := args["directory"].(string)
	if !ok {
		return nil, fmt.Errorf("directory parameter required")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf(
					"Duplicate groups for: %s\n\nThis would return smart groups of duplicate files with:\n- Exact matches\n- Similar files\n- Potential space savings\n- Recommended files to keep/remove",
					directory,
				),
			},
		},
	}, nil
}

func checkFileHashTool(args map[string]interface{}) (interface{}, error) {
	filePath, ok := args["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path parameter required")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf(
					"File: %s\n\nThis would return:\n- Content hash (SHA256)\n- Merkle root\n- File size\n- Chunk count",
					filePath,
				),
			},
		},
	}, nil
}
