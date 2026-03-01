// Command services-voice is the entry point for the services.voice plugin
// binary. It provides 8 MCP tools for text-to-speech and speech-to-text.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orchestra-mcp/plugin-services-voice/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

func main() {
	builder := plugin.New("services.voice").
		Version("0.1.0").
		Description("Text-to-speech and speech-to-text voice services").
		Author("Orchestra").
		Binary("services-voice")

	tp := &internal.VoicePlugin{}
	tp.RegisterTools(builder)

	p := builder.BuildWithTools()
	p.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Run(ctx); err != nil {
		log.Fatalf("services.voice: %v", err)
	}
}
