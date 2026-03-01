package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// SttListenSchema returns the JSON Schema for the stt_listen tool.
func SttListenSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"duration": map[string]any{
				"type":        "integer",
				"description": "Duration in seconds to listen (default: 5)",
			},
		},
	})
	return s
}

// SttListen returns a handler that explains the microphone STT requirements.
func SttListen() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		return helpers.TextResult(
			"STT via microphone requires additional setup. Install SpeechRecognition and configure OPENAI_API_KEY for Whisper.",
		), nil
	}
}
