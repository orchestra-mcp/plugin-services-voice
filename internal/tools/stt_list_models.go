package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// sttModels maps provider names to their available STT models.
var sttModels = map[string][]string{
	"system": {
		"macos:SFSpeechRecognizer (built-in, requires macOS 10.15+)",
		"linux:Vosk (requires vosk-model installed)",
	},
	"openai": {
		"whisper-1 (general-purpose, multilingual)",
	},
	"google": {
		"latest_long (long audio, high accuracy)",
		"latest_short (short audio, low latency)",
		"command_and_search",
	},
	"elevenlabs": {
		"scribe_v1 (multilingual)",
	},
}

// SttListModelsSchema returns the JSON Schema for the stt_list_models tool.
func SttListModelsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"provider": map[string]any{
				"type":        "string",
				"description": "Filter by provider: system, openai, google, elevenlabs (omit for all)",
			},
		},
	})
	return s
}

// SttListModels returns a handler that lists available STT models, optionally
// filtered by provider.
func SttListModels() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		provider := helpers.GetString(req.Arguments, "provider")

		var sb strings.Builder
		sb.WriteString("## Available STT Models\n\n")

		if provider != "" {
			models, ok := sttModels[provider]
			if !ok {
				return helpers.ErrorResult("validation_error",
					fmt.Sprintf("unknown provider %q: must be one of system, openai, google, elevenlabs", provider)), nil
			}
			sb.WriteString(fmt.Sprintf("### %s\n", provider))
			for _, m := range models {
				sb.WriteString(fmt.Sprintf("- %s\n", m))
			}
		} else {
			for _, p := range []string{"system", "openai", "google", "elevenlabs"} {
				sb.WriteString(fmt.Sprintf("### %s\n", p))
				for _, m := range sttModels[p] {
					sb.WriteString(fmt.Sprintf("- %s\n", m))
				}
				sb.WriteString("\n")
			}
		}

		return helpers.TextResult(sb.String()), nil
	}
}
