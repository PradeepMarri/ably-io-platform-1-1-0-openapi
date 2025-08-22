package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/platform-api/mcp-server/config"
	"github.com/platform-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetmetadataofchannelHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		channel_idVal, ok := args["channel_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: channel_id"), nil
		}
		channel_id, ok := channel_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: channel_id"), nil
		}
		url := fmt.Sprintf("%s/channels/%s", cfg.BaseURL, channel_id)
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
		var result models.ChannelDetails
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

func CreateGetmetadataofchannelTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_channels_channel_id",
		mcp.WithDescription("Get metadata of a channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("The [Channel's ID](https://www.ably.io/documentation/rest/channels).")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetmetadataofchannelHandler(cfg),
	}
}
