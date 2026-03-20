package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/globaldb"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-services-voice/internal/tts"
	"google.golang.org/protobuf/types/known/structpb"
)

// TtsSpeakSchema returns the JSON Schema for the tts_speak tool.
func TtsSpeakSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]any{
				"type":        "string",
				"description": "The text to speak aloud",
			},
			"voice": map[string]any{
				"type":        "string",
				"description": "Optional voice name (e.g. 'Samantha' on macOS). Falls back to notify.voice_name config.",
			},
			"speed": map[string]any{
				"type":        "string",
				"description": "Optional speed in words per minute (e.g. '180'). Falls back to notify.voice_speed config.",
			},
			"volume": map[string]any{
				"type":        "string",
				"description": "Optional volume from 0.0 to 1.0 (e.g. '0.8'). Falls back to notify.voice_volume config. Note: macOS say does not support volume.",
			},
		},
		"required": []any{"text"},
	})
	return s
}

// TtsSpeak returns a handler that converts text to speech using the system TTS engine.
func TtsSpeak() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "text"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		text := helpers.GetString(req.Arguments, "text")
		voice := helpers.GetString(req.Arguments, "voice")
		speed := helpers.GetString(req.Arguments, "speed")
		volume := helpers.GetString(req.Arguments, "volume")

		// Fall back to stored preferences, then platform defaults.
		if voice == "" {
			voice = globaldb.GetConfig("notify.voice_name")
		}
		// If no voice set, use platform system default (Siri on macOS).
		// Users can override via notify_config settings.
		if speed == "" {
			speed = globaldb.GetConfig("notify.voice_speed")
		}
		if volume == "" {
			volume = globaldb.GetConfig("notify.voice_volume")
		}

		_, err := tts.Speak(ctx, text, voice, speed, volume)
		if err != nil {
			return helpers.ErrorResult("tts_error", fmt.Sprintf("TTS failed: %v", err)), nil
		}

		msg := fmt.Sprintf("Spoke: %q", text)
		if voice != "" {
			msg = fmt.Sprintf("Spoke: %q (voice: %s)", text, voice)
		}
		return helpers.TextResult(msg), nil
	}
}
