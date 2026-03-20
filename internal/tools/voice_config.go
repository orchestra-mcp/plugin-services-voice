package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/globaldb"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// VoiceConfigSchema returns the JSON Schema for the voice_config tool.
func VoiceConfigSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": `Action to perform: "get" to read current voice config, "set" to update voice config`,
				"enum":        []any{"get", "set"},
			},
			"default_voice": map[string]any{
				"type":        "string",
				"description": "Set the default voice name (e.g. 'Samantha' on macOS)",
			},
			"speed": map[string]any{
				"type":        "string",
				"description": "TTS speed in words per minute (e.g. '180')",
			},
			"volume": map[string]any{
				"type":        "string",
				"description": "TTS volume from 0.0 to 1.0 (e.g. '0.8')",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Set the default language code (e.g. 'en-US')",
			},
		},
		"required": []any{"action"},
	})
	return s
}

// VoiceConfig returns a handler that gets or sets voice configuration via globaldb.
func VoiceConfig() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		action := helpers.GetString(req.Arguments, "action")
		if action == "" {
			// Auto-detect: if any set-fields are provided, treat as "set".
			if helpers.GetString(req.Arguments, "default_voice") != "" ||
				helpers.GetString(req.Arguments, "speed") != "" ||
				helpers.GetString(req.Arguments, "volume") != "" ||
				helpers.GetString(req.Arguments, "language") != "" {
				action = "set"
			} else {
				action = "get"
			}
		}

		switch action {
		case "get":
			voiceName := globaldb.GetConfig("notify.voice_name")
			speed := globaldb.GetConfig("notify.voice_speed")
			volume := globaldb.GetConfig("notify.voice_volume")
			language := globaldb.GetConfig("notify.voice_language")

			if voiceName == "" {
				voiceName = "(system default)"
			}
			if speed == "" {
				speed = "(default)"
			}
			if volume == "" {
				volume = "(default)"
			}
			if language == "" {
				language = "(system default)"
			}

			result, err := helpers.JSONResult(map[string]any{
				"voice_name": voiceName,
				"speed":      speed,
				"volume":     volume,
				"language":   language,
			})
			if err != nil {
				return helpers.TextResult(fmt.Sprintf("voice_name=%s, speed=%s, volume=%s, language=%s",
					voiceName, speed, volume, language)), nil
			}
			return result, nil

		case "set":
			var parts []string
			if v := helpers.GetString(req.Arguments, "default_voice"); v != "" {
				_ = globaldb.SetConfig("notify.voice_name", v)
				parts = append(parts, fmt.Sprintf("voice_name=%s", v))
			}
			if v := helpers.GetString(req.Arguments, "speed"); v != "" {
				_ = globaldb.SetConfig("notify.voice_speed", v)
				parts = append(parts, fmt.Sprintf("speed=%s", v))
			}
			if v := helpers.GetString(req.Arguments, "volume"); v != "" {
				_ = globaldb.SetConfig("notify.voice_volume", v)
				parts = append(parts, fmt.Sprintf("volume=%s", v))
			}
			if v := helpers.GetString(req.Arguments, "language"); v != "" {
				_ = globaldb.SetConfig("notify.voice_language", v)
				parts = append(parts, fmt.Sprintf("language=%s", v))
			}
			if len(parts) == 0 {
				return helpers.ErrorResult("validation_error", "at least one config field must be provided for action=set"), nil
			}
			return helpers.TextResult(fmt.Sprintf("Voice config updated: %s", strings.Join(parts, ", "))), nil

		default:
			return helpers.ErrorResult("validation_error", fmt.Sprintf("unknown action %q: must be 'get' or 'set'", action)), nil
		}
	}
}
