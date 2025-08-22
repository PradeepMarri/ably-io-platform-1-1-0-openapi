package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/platform-api/mcp-server/config"
	"github.com/platform-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func PublishmessagestochannelHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		// Create properly typed request body using the generated schema
		var requestBody models.Message
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/channels/%s/messages", cfg.BaseURL, channel_id)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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

func CreatePublishmessagestochannelTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_channels_channel_id_messages",
		mcp.WithDescription("Publish a message to a channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("The [Channel's ID](https://www.ably.io/documentation/rest/channels).")),
		mcp.WithNumber("timestamp", mcp.Description("Input parameter: Timestamp when the message was received by the Ably, as milliseconds since the epoch.")),
		mcp.WithString("clientId", mcp.Description("Input parameter: The [client ID](https://www.ably.io/documentation/core-features/authentication#identified-clients) of the publisher of this message.")),
		mcp.WithString("connectionId", mcp.Description("Input parameter: The connection ID of the publisher of this message.")),
		mcp.WithString("data", mcp.Description("Input parameter: The string encoded payload, with the encoding specified below.")),
		mcp.WithString("encoding", mcp.Description("Input parameter: This will typically be empty as all messages received from Ably are automatically decoded client-side using this value. However, if the message encoding cannot be processed, this attribute will contain the remaining transformations not applied to the data payload.")),
		mcp.WithObject("extras", mcp.Description("Input parameter: Extras object. Currently only allows for [push](https://www.ably.io/documentation/general/push/publish#channel-broadcast-example) extra.")),
		mcp.WithString("id", mcp.Description("Input parameter: A Unique ID that can be specified by the publisher for [idempotent publishing](https://www.ably.io/documentation/rest/messages#idempotent).")),
		mcp.WithString("name", mcp.Description("Input parameter: The event name, if provided.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    PublishmessagestochannelHandler(cfg),
	}
}
