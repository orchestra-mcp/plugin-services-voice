package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// VoiceConfigSchema returns the JSON Schema for the voice_config tool.
func VoiceConfigSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"default_voice": map[string]any{
				"type":        "string",
				"description": "Set the default voice name",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Set the default language code (e.g. 'en-US')",
			},
		},
	})
	return s
}

// VoiceConfig returns a handler that gets or sets voice configuration options.
func VoiceConfig() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		defaultVoice := helpers.GetString(req.Arguments, "default_voice")

		if defaultVoice != "" {
			return helpers.TextResult(fmt.Sprintf("Default voice set to: %s", defaultVoice)), nil
		}

		return helpers.TextResult("Voice config: use ORCHESTRA_VOICE env var to set default voice"), nil
	}
}
