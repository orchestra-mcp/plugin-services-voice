package servicesvoice

import (
	"github.com/orchestra-mcp/plugin-services-voice/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all voice tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	vp := &internal.VoicePlugin{}
	vp.RegisterTools(builder)
}
