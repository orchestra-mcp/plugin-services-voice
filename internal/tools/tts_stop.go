package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// TtsStopSchema returns the JSON Schema for the tts_stop tool.
func TtsStopSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	})
	return s
}

// TtsStop returns a handler that kills any currently running TTS process.
func TtsStop() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		var target string
		switch runtime.GOOS {
		case "darwin":
			target = "say"
		default:
			target = "espeak"
		}

		cmd := exec.CommandContext(ctx, "pkill", target)
		out, err := cmd.CombinedOutput()
		if err != nil {
			// pkill exits with code 1 when no matching process is found — treat
			// this as a non-fatal condition and still report success.
			_ = out
			return helpers.TextResult(fmt.Sprintf("TTS stopped (no active %s process found)", target)), nil
		}
		return helpers.TextResult("TTS stopped"), nil
	}
}
