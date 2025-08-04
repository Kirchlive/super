// Package main is the entry point for the hello plugin
package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/opencode-superclaude/examples/simple-plugin/shared"
)

func main() {
	// Create an instance of our plugin
	helloPlugin := &HelloPlugin{}

	// Serve the plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"command": &shared.CommandPluginImpl{Impl: helloPlugin},
		},
		
		// A non-nil value here enables gRPC serving for this plugin
		GRPCServer: plugin.DefaultGRPCServer,
	})
}