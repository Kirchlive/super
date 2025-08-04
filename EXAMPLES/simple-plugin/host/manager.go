// Package main implements the plugin manager for the host application
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/opencode-superclaude/examples/simple-plugin/shared"
)

// PluginInfo stores metadata about a loaded plugin
type PluginInfo struct {
	Name         string
	Version      string
	Path         string
	Capabilities []string
	Client       *plugin.Client
	Instance     shared.CommandPlugin
}

// PluginManager manages the lifecycle of plugins
type PluginManager struct {
	plugins map[string]*PluginInfo
	mu      sync.RWMutex
}

// NewPluginManager creates a new plugin manager instance
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*PluginInfo),
	}
}

// DiscoverPlugins searches for and loads plugins from the specified directory
func (pm *PluginManager) DiscoverPlugins(dir string) error {
	log.Printf("Discovering plugins in: %s", dir)
	
	// Ensure plugin directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}
	
	// Find all plugin binaries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		pluginPath := filepath.Join(dir, entry.Name())
		log.Printf("Found potential plugin: %s", pluginPath)
		
		// Load the plugin
		if err := pm.LoadPlugin(pluginPath); err != nil {
			log.Printf("Failed to load plugin %s: %v", pluginPath, err)
			continue
		}
	}
	
	return nil
}

// LoadPlugin loads a single plugin from the specified path
func (pm *PluginManager) LoadPlugin(path string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// Create plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})
	
	// Connect to the plugin
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to create RPC client: %w", err)
	}
	
	// Get the plugin instance
	raw, err := rpcClient.Dispense("command")
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin: %w", err)
	}
	
	// Cast to our interface
	pluginInstance, ok := raw.(shared.CommandPlugin)
	if !ok {
		client.Kill()
		return fmt.Errorf("plugin does not implement CommandPlugin interface")
	}
	
	// Get plugin metadata
	name := pluginInstance.Name()
	version := pluginInstance.Version()
	capabilities := pluginInstance.GetCapabilities()
	
	// Store plugin info
	info := &PluginInfo{
		Name:         name,
		Version:      version,
		Path:         path,
		Capabilities: capabilities,
		Client:       client,
		Instance:     pluginInstance,
	}
	
	pm.plugins[name] = info
	log.Printf("Loaded plugin: %s v%s", name, version)
	
	return nil
}

// ExecutePlugin executes a command on the specified plugin
func (pm *PluginManager) ExecutePlugin(name string, args map[string]interface{}) (string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	info, exists := pm.plugins[name]
	if !exists {
		return "", fmt.Errorf("plugin not found: %s", name)
	}
	
	// Execute the plugin
	result, err := info.Instance.Execute(args)
	if err != nil {
		return "", fmt.Errorf("plugin execution failed: %w", err)
	}
	
	return result, nil
}

// ListPlugins returns information about all loaded plugins
func (pm *PluginManager) ListPlugins() []PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var plugins []PluginInfo
	for _, info := range pm.plugins {
		plugins = append(plugins, *info)
	}
	
	return plugins
}

// UnloadPlugin unloads a specific plugin
func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	
	// Kill the plugin process
	info.Client.Kill()
	
	// Remove from registry
	delete(pm.plugins, name)
	log.Printf("Unloaded plugin: %s", name)
	
	return nil
}

// ReloadPlugin reloads a plugin (useful for hot-reload)
func (pm *PluginManager) ReloadPlugin(name string) error {
	pm.mu.RLock()
	info, exists := pm.plugins[name]
	pm.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	
	path := info.Path
	
	// Unload the current version
	if err := pm.UnloadPlugin(name); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}
	
	// Load the new version
	if err := pm.LoadPlugin(path); err != nil {
		return fmt.Errorf("failed to reload plugin: %w", err)
	}
	
	log.Printf("Reloaded plugin: %s", name)
	return nil
}

// Shutdown gracefully shuts down all plugins
func (pm *PluginManager) Shutdown() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	for name, info := range pm.plugins {
		log.Printf("Shutting down plugin: %s", name)
		info.Client.Kill()
	}
	
	pm.plugins = make(map[string]*PluginInfo)
}