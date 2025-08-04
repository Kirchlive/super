// Package main is the host application that loads and manages plugins
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Set up logging
	log.SetPrefix("[HOST] ")
	log.SetFlags(log.Ltime | log.Lshortfile)
	
	fmt.Println("=== OpenCode Plugin System Demo ===")
	fmt.Println()
	
	// Create plugin manager
	manager := NewPluginManager()
	
	// Discover and load plugins
	log.Println("Starting plugin system...")
	if err := manager.DiscoverPlugins("./plugins"); err != nil {
		log.Fatalf("Failed to discover plugins: %v", err)
	}
	
	// List loaded plugins
	plugins := manager.ListPlugins()
	fmt.Printf("Loaded %d plugin(s):\n", len(plugins))
	for _, p := range plugins {
		fmt.Printf("  - %s (v%s)\n", p.Name, p.Version)
		fmt.Printf("    Capabilities: %v\n", p.Capabilities)
	}
	fmt.Println()
	
	// Execute plugin command
	pluginName := "hello"
	log.Printf("Executing command on plugin: %s", pluginName)
	
	// Test different greeting types
	testCases := []struct {
		args map[string]interface{}
		desc string
	}{
		{
			args: map[string]interface{}{
				"name": "Developer",
				"type": "standard",
			},
			desc: "Standard greeting",
		},
		{
			args: map[string]interface{}{
				"name": "Enterprise User",
				"type": "formal",
			},
			desc: "Formal greeting",
		},
		{
			args: map[string]interface{}{
				"name": "Coder",
				"type": "casual",
			},
			desc: "Casual greeting",
		},
		{
			args: map[string]interface{}{
				"name": "System",
				"type": "technical",
			},
			desc: "Technical greeting",
		},
	}
	
	for _, tc := range testCases {
		fmt.Printf("\n%s:\n", tc.desc)
		result, err := manager.ExecutePlugin(pluginName, tc.args)
		if err != nil {
			log.Printf("Error executing plugin: %v", err)
			continue
		}
		fmt.Printf("  Response: %s\n", result)
	}
	
	// Demonstrate plugin hot-reload capability
	fmt.Println("\n--- Hot Reload Demo ---")
	fmt.Println("In a real implementation, you could:")
	fmt.Println("1. Modify the plugin source")
	fmt.Println("2. Rebuild the plugin")
	fmt.Println("3. The host would detect changes and reload")
	fmt.Println()
	
	// Clean shutdown
	log.Println("Shutting down plugin system...")
	manager.Shutdown()
	
	fmt.Println("\n=== Demo Complete ===")
	os.Exit(0)
}