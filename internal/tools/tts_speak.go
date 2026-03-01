package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
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
				"description": "Optional voice name (e.g. 'Samantha' on macOS)",
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

		_, err := tts.Speak(ctx, text, voice)
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
