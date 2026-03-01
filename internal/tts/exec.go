package tts

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Speak invokes the platform TTS engine to speak the given text.
// On macOS it uses the built-in `say` command; on all other platforms it uses
// `espeak`. An optional voice name may be supplied.
func Speak(ctx context.Context, text, voice string) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		args := []string{"say"}
		if voice != "" {
			args = append(args, "-v", voice)
		}
		args = append(args, text)
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	default:
		args := []string{"espeak", text}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListVoices returns the list of available voices from the platform TTS engine.
// On macOS it runs `say -v ?`; on all other platforms it runs `espeak --voices`.
func ListVoices(ctx context.Context) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "say", "-v", "?")
	default:
		cmd = exec.CommandContext(ctx, "espeak", "--voices")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
