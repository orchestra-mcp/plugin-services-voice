package internal

import (
	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-services-voice/internal/tools"
)

// VoicePlugin registers all TTS and STT tools with the plugin builder.
type VoicePlugin struct{}

// RegisterTools registers all 8 voice tools with the plugin builder.
func (vp *VoicePlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("tts_speak",
		"Convert text to speech using the system TTS engine",
		tools.TtsSpeakSchema(), tools.TtsSpeak())

	builder.RegisterTool("tts_speak_provider",
		"Convert text to speech using an external provider (ElevenLabs, OpenAI, Google)",
		tools.TtsSpeakProviderSchema(), tools.TtsSpeakProvider())

	builder.RegisterTool("tts_list_voices",
		"List available TTS voices on the current platform",
		tools.TtsListVoicesSchema(), tools.TtsListVoices())

	builder.RegisterTool("tts_stop",
		"Stop any currently playing TTS audio",
		tools.TtsStopSchema(), tools.TtsStop())

	builder.RegisterTool("stt_listen",
		"Listen via microphone and transcribe speech to text",
		tools.SttListenSchema(), tools.SttListen())

	builder.RegisterTool("stt_transcribe_file",
		"Transcribe an audio file to text",
		tools.SttTranscribeFileSchema(), tools.SttTranscribeFile())

	builder.RegisterTool("stt_list_models",
		"List available STT models, optionally filtered by provider",
		tools.SttListModelsSchema(), tools.SttListModels())

	builder.RegisterTool("voice_config",
		"Get or set voice configuration options",
		tools.VoiceConfigSchema(), tools.VoiceConfig())
}
