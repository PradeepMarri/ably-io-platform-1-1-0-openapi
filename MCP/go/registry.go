package main

import (
	"github.com/platform-api/mcp-server/config"
	"github.com/platform-api/mcp-server/models"
	tools_stats "github.com/platform-api/mcp-server/tools/stats"
	tools_status "github.com/platform-api/mcp-server/tools/status"
	tools_push "github.com/platform-api/mcp-server/tools/push"
	tools_history "github.com/platform-api/mcp-server/tools/history"
	tools_publishing "github.com/platform-api/mcp-server/tools/publishing"
	tools_authentication "github.com/platform-api/mcp-server/tools/authentication"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_stats.CreateGettimeTool(cfg),
		tools_status.CreateGetmetadataofallchannelsTool(cfg),
		tools_status.CreateGetpresenceofchannelTool(cfg),
		tools_push.CreateUnregisterallpushdevicesTool(cfg),
		tools_push.CreateGetregisteredpushdevicesTool(cfg),
		tools_push.CreateRegisterpushdeviceTool(cfg),
		tools_stats.CreateGetstatsTool(cfg),
		tools_history.CreateGetmessagesbychannelTool(cfg),
		tools_publishing.CreatePublishmessagestochannelTool(cfg),
		tools_push.CreateSubscribepushdevicetochannelTool(cfg),
		tools_push.CreateDeletepushdevicedetailsTool(cfg),
		tools_push.CreateGetpushsubscriptionsonchannelsTool(cfg),
		tools_push.CreateGetchannelswithpushsubscribersTool(cfg),
		tools_push.CreatePublishpushnotificationtodevicesTool(cfg),
		tools_status.CreateGetmetadataofchannelTool(cfg),
		tools_history.CreateGetpresencehistoryofchannelTool(cfg),
		tools_authentication.CreateRequestaccesstokenTool(cfg),
		tools_push.CreateUnregisterpushdeviceTool(cfg),
		tools_push.CreateGetpushdevicedetailsTool(cfg),
		tools_push.CreatePatchpushdevicedetailsTool(cfg),
		tools_push.CreatePutpushdevicedetailsTool(cfg),
		tools_push.CreateUpdatepushdevicedetailsTool(cfg),
	}
}
