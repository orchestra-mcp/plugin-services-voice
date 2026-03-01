package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-services-voice/internal/tts"
	"google.golang.org/protobuf/types/known/structpb"
)

// TtsListVoicesSchema returns the JSON Schema for the tts_list_voices tool.
func TtsListVoicesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	})
	return s
}

// TtsListVoices returns a handler that lists available TTS voices on the current platform.
func TtsListVoices() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		result, err := tts.ListVoices(ctx)
		if err != nil {
			return helpers.ErrorResult("tts_error", fmt.Sprintf("Failed to list voices: %v", err)), nil
		}
		return helpers.TextResult(result), nil
	}
}
