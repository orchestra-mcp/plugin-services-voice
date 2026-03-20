package tts

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Speak invokes the platform TTS engine to speak the given text.
// On macOS it uses the built-in `say` command; on all other platforms it uses
// `espeak`. Optional voice, speed (words per minute), and volume (0.0-1.0) may be supplied.
// Note: macOS `say` does not support volume control natively — only speed via -r flag.
func Speak(ctx context.Context, text, voice, speed, volume string) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		args := []string{}
		if voice != "" {
			args = append(args, "-v", voice)
		}
		if speed != "" {
			args = append(args, "-r", speed)
		}
		args = append(args, text)
		cmd = exec.CommandContext(ctx, "say", args...)
	default:
		args := []string{}
		if speed != "" {
			args = append(args, "-s", speed)
		}
		if volume != "" {
			if vol, err := strconv.ParseFloat(volume, 64); err == nil {
				// espeak amplitude is 0-200, map 0.0-1.0 to 0-200
				amp := int(vol * 200)
				if amp < 0 {
					amp = 0
				}
				if amp > 200 {
					amp = 200
				}
				args = append(args, "-a", strconv.Itoa(amp))
			}
		}
		args = append(args, text)
		cmd = exec.CommandContext(ctx, "espeak", args...)
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
