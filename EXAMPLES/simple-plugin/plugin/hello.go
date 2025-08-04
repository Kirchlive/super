// Package main implements the hello plugin functionality
package main

import (
	"fmt"
	"log"
)

// HelloPlugin is a simple plugin that demonstrates the plugin architecture
type HelloPlugin struct{}

// Name returns the plugin's unique identifier
func (p *HelloPlugin) Name() string {
	return "hello"
}

// Version returns the plugin's version
func (p *HelloPlugin) Version() string {
	return "1.0.0"
}

// Execute runs the plugin's main functionality
func (p *HelloPlugin) Execute(args map[string]interface{}) (string, error) {
	log.Println("[PLUGIN] Executing hello command")
	
	// Extract name from args, with default
	name := "World"
	if n, ok := args["name"].(string); ok && n != "" {
		name = n
	}
	
	// Extract greeting type
	greetingType := "standard"
	if t, ok := args["type"].(string); ok {
		greetingType = t
	}
	
	// Generate response based on type
	var response string
	switch greetingType {
	case "formal":
		response = fmt.Sprintf("Greetings, %s. Welcome to the SuperClaude integration platform.", name)
	case "casual":
		response = fmt.Sprintf("Hey %s! Ready to enhance OpenCode with AI?", name)
	case "technical":
		response = fmt.Sprintf("Plugin 'hello' v%s initialized. Target: %s. Integration: operational.", p.Version(), name)
	default:
		response = fmt.Sprintf("Hello %s from SuperClaude integration!", name)
	}
	
	log.Printf("[PLUGIN] Generated response: %s", response)
	return response, nil
}

// GetCapabilities returns a list of capabilities this plugin provides
func (p *HelloPlugin) GetCapabilities() []string {
	return []string{
		"greet",           // Basic greeting
		"greet.formal",    // Formal greeting
		"greet.casual",    // Casual greeting
		"greet.technical", // Technical greeting
		"plugin.info",     // Plugin information
	}
}

// Additional methods that could be added in a real plugin:

// Initialize would set up any resources the plugin needs
func (p *HelloPlugin) Initialize(config map[string]interface{}) error {
	log.Println("[PLUGIN] Hello plugin initialized")
	// In a real plugin, this might:
	// - Connect to databases
	// - Load configuration
	// - Set up event listeners
	return nil
}

// Cleanup would clean up resources when the plugin shuts down
func (p *HelloPlugin) Cleanup() error {
	log.Println("[PLUGIN] Hello plugin shutting down")
	// In a real plugin, this might:
	// - Close database connections
	// - Save state
	// - Clean up temporary files
	return nil
}