package tools

import (
	"context"
	"fmt"
	"os"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// SttTranscribeFileSchema returns the JSON Schema for the stt_transcribe_file tool.
func SttTranscribeFileSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Path to the audio file to transcribe",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Optional BCP-47 language code (e.g. 'en', 'fr')",
			},
		},
		"required": []any{"file_path"},
	})
	return s
}

// SttTranscribeFile returns a handler that checks the given file and explains
// the transcription requirements.
func SttTranscribeFile() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "file_path"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		filePath := helpers.GetString(req.Arguments, "file_path")

		info, err := os.Stat(filePath)
		if err != nil {
			return helpers.ErrorResult("file_not_found", fmt.Sprintf("File not found: %s", filePath)), nil
		}

		return helpers.TextResult(
			fmt.Sprintf(
				"Transcription requires OPENAI_API_KEY. Set the env var and retry. File: %s (%d bytes)",
				filePath, info.Size(),
			),
		), nil
	}
}
