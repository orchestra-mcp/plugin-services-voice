package tools

import (
	"context"
	"os"
	"strings"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- helpers ----------

func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

func getText(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	if r := resp.GetResult(); r != nil {
		if f := r.GetFields(); f != nil {
			if tf, ok := f["text"]; ok {
				return tf.GetStringValue()
			}
		}
	}
	return ""
}

// ---------- tts_speak ----------

func TestTtsSpeak_MissingText(t *testing.T) {
	resp := callTool(t, TtsSpeak(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing text")
	}
}

func TestTtsSpeak_ValidText(t *testing.T) {
	// May fail if TTS not available — that is fine (tts_error is acceptable).
	resp := callTool(t, TtsSpeak(), map[string]any{"text": "hello"})
	_ = resp
}

// ---------- tts_speak_provider ----------

func TestTtsSpeakProvider_MissingText(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{"provider": "openai"})
	if !isError(resp) {
		t.Error("expected validation_error for missing text")
	}
}

func TestTtsSpeakProvider_MissingProvider(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{"text": "hello"})
	if !isError(resp) {
		t.Error("expected validation_error for missing provider")
	}
}

func TestTtsSpeakProvider_InvalidProvider(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{
		"text":     "hello",
		"provider": "amazon",
	})
	if !isError(resp) {
		t.Error("expected validation_error for unknown provider")
	}
}

func TestTtsSpeakProvider_SystemProvider(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{
		"text":     "hello world",
		"provider": "system",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "system") {
		t.Errorf("expected 'system' in response, got: %s", txt)
	}
}

func TestTtsSpeakProvider_OpenAIProvider(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{
		"text":     "hello",
		"provider": "openai",
		"model":    "tts-1",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "OPENAI_API_KEY") {
		t.Errorf("expected OPENAI_API_KEY in response, got: %s", txt)
	}
}

func TestTtsSpeakProvider_ElevenLabsProvider(t *testing.T) {
	resp := callTool(t, TtsSpeakProvider(), map[string]any{
		"text":     "hello",
		"provider": "elevenlabs",
		"voice":    "Rachel",
	})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "ELEVENLABS_API_KEY") {
		t.Errorf("expected ELEVENLABS_API_KEY in response, got: %s", txt)
	}
}

// ---------- tts_list_voices ----------

func TestTtsListVoices_NoArgs(t *testing.T) {
	// May fail if `say`/`espeak` not available — acceptable.
	resp := callTool(t, TtsListVoices(), map[string]any{})
	_ = resp
}

// ---------- tts_stop ----------

func TestTtsStop_NoArgs(t *testing.T) {
	// pkill may not find a process — always returns success.
	resp := callTool(t, TtsStop(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
}

// ---------- stt_listen ----------

func TestSttListen_NoArgs(t *testing.T) {
	resp := callTool(t, SttListen(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
}

// ---------- stt_transcribe_file ----------

func TestSttTranscribeFile_MissingFilePath(t *testing.T) {
	resp := callTool(t, SttTranscribeFile(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing file_path")
	}
}

func TestSttTranscribeFile_NonexistentFile(t *testing.T) {
	resp := callTool(t, SttTranscribeFile(), map[string]any{
		"file_path": "/tmp/no-such-audio-orchestra-xyz.wav",
	})
	if !isError(resp) {
		t.Error("expected file_not_found for nonexistent file")
	}
}

func TestSttTranscribeFile_ExistingFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "audio-*.wav")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	resp := callTool(t, SttTranscribeFile(), map[string]any{"file_path": f.Name()})
	if isError(resp) {
		t.Errorf("unexpected error for existing file: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "OPENAI_API_KEY") {
		t.Errorf("expected OPENAI_API_KEY hint in response, got: %s", txt)
	}
}

// ---------- stt_list_models ----------

func TestSttListModels_NoArgs(t *testing.T) {
	resp := callTool(t, SttListModels(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	for _, provider := range []string{"system", "openai", "google", "elevenlabs"} {
		if !strings.Contains(txt, provider) {
			t.Errorf("expected provider %q in response", provider)
		}
	}
}

func TestSttListModels_FilterByProvider(t *testing.T) {
	resp := callTool(t, SttListModels(), map[string]any{"provider": "openai"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "whisper") {
		t.Errorf("expected whisper model in openai response, got: %s", txt)
	}
}

func TestSttListModels_InvalidProvider(t *testing.T) {
	resp := callTool(t, SttListModels(), map[string]any{"provider": "amazon"})
	if !isError(resp) {
		t.Error("expected validation_error for unknown provider")
	}
}

// ---------- voice_config ----------

func TestVoiceConfig_NoArgs(t *testing.T) {
	resp := callTool(t, VoiceConfig(), map[string]any{})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
}

func TestVoiceConfig_SetDefaultVoice(t *testing.T) {
	resp := callTool(t, VoiceConfig(), map[string]any{"default_voice": "Samantha"})
	if isError(resp) {
		t.Errorf("unexpected error: %s", getText(resp))
	}
	txt := getText(resp)
	if !strings.Contains(txt, "Samantha") {
		t.Errorf("expected voice name in response, got: %s", txt)
	}
}
