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

func GetpushdevicedetailsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		device_idVal, ok := args["device_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: device_id"), nil
		}
		device_id, ok := device_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: device_id"), nil
		}
		url := fmt.Sprintf("%s/push/deviceRegistrations/%s", cfg.BaseURL, device_id)
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

func CreateGetpushdevicedetailsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_push_deviceRegistrations_device_id",
		mcp.WithDescription("Get a device registration"),
		mcp.WithString("device_id", mcp.Required(), mcp.Description("Device's ID.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetpushdevicedetailsHandler(cfg),
	}
}
