package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/platform-api/mcp-server/config"
	"github.com/platform-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetmetadataofallchannelsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["limit"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("limit=%v", val))
		}
		if val, ok := args["prefix"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("prefix=%v", val))
		}
		if val, ok := args["by"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("by=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/channels%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		if cfg.BasicAuth != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", cfg.BasicAuth))
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.Error
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetmetadataofallchannelsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_channels",
		mcp.WithDescription("Enumerate all active channels of the application"),
		mcp.WithNumber("limit", mcp.Description("")),
		mcp.WithString("prefix", mcp.Description("Optionally limits the query to only those channels whose name starts with the given prefix")),
		mcp.WithString("by", mcp.Description("optionally specifies whether to return just channel names (by=id) or ChannelDetails (by=value)")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetmetadataofallchannelsHandler(cfg),
	}
}
