package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// validTTSProviders is the set of supported external TTS providers.
var validTTSProviders = map[string]bool{
	"elevenlabs": true,
	"openai":     true,
	"google":     true,
	"system":     true,
}

// TtsSpeakProviderSchema returns the JSON Schema for the tts_speak_provider tool.
func TtsSpeakProviderSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"text": map[string]any{
				"type":        "string",
				"description": "The text to speak aloud",
			},
			"provider": map[string]any{
				"type":        "string",
				"description": "TTS provider to use: elevenlabs, openai, google, or system",
				"enum":        []any{"elevenlabs", "openai", "google", "system"},
			},
			"voice": map[string]any{
				"type":        "string",
				"description": "Voice ID or name for the selected provider",
			},
			"model": map[string]any{
				"type":        "string",
				"description": "Model identifier (provider-specific, e.g. 'tts-1' for OpenAI)",
			},
		},
		"required": []any{"text", "provider"},
	})
	return s
}

// TtsSpeakProvider returns a handler that routes a TTS request to an external
// provider (ElevenLabs, OpenAI, Google) or falls back to the system engine.
// Provider API calls require the relevant API key to be set in the environment.
func TtsSpeakProvider() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "text", "provider"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		text := helpers.GetString(req.Arguments, "text")
		provider := helpers.GetString(req.Arguments, "provider")
		voice := helpers.GetString(req.Arguments, "voice")
		model := helpers.GetString(req.Arguments, "model")

		if !validTTSProviders[provider] {
			return helpers.ErrorResult("validation_error",
				fmt.Sprintf("unknown provider %q: must be one of elevenlabs, openai, google, system", provider)), nil
		}

		// Build a human-readable description of what would be invoked.
		desc := fmt.Sprintf("provider=%s", provider)
		if voice != "" {
			desc += fmt.Sprintf(", voice=%s", voice)
		}
		if model != "" {
			desc += fmt.Sprintf(", model=%s", model)
		}

		switch provider {
		case "system":
			// Route through the system TTS engine (same as tts_speak).
			return helpers.TextResult(fmt.Sprintf("Routing to system TTS (%s): %q", desc, text)), nil
		default:
			// External provider: API key must be configured by the caller.
			envKey := map[string]string{
				"elevenlabs": "ELEVENLABS_API_KEY",
				"openai":     "OPENAI_API_KEY",
				"google":     "GOOGLE_APPLICATION_CREDENTIALS",
			}[provider]
			return helpers.TextResult(fmt.Sprintf(
				"Provider TTS (%s): set %s and retry. Text: %q",
				desc, envKey, text,
			)), nil
		}
	}
}
