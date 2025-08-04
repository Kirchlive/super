// Package shared contains interfaces and types shared between host and plugins
package shared

import (
	"github.com/hashicorp/go-plugin"
)

// CommandPlugin is the interface that all OpenCode plugins must implement
type CommandPlugin interface {
	// Name returns the plugin's unique identifier
	Name() string
	
	// Version returns the plugin's version
	Version() string
	
	// Execute runs the plugin's main functionality
	Execute(args map[string]interface{}) (string, error)
	
	// GetCapabilities returns a list of capabilities this plugin provides
	GetCapabilities() []string
}

// CommandPluginRPC is the RPC implementation of CommandPlugin
type CommandPluginRPC struct {
	client *plugin.Client
	broker *plugin.MuxBroker
}

// Handshake is a common handshake that is shared by plugin and host
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "OPENCODE_PLUGIN",
	MagicCookieValue: "superclaude",
}

// PluginMap is the map of plugins we can dispense
var PluginMap = map[string]plugin.Plugin{
	"command": &CommandPluginImpl{},
}

// CommandPluginImpl is the implementation of plugin.Plugin for CommandPlugin
type CommandPluginImpl struct {
	Impl CommandPlugin
}

func (p *CommandPluginImpl) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &CommandPluginRPCServer{Impl: p.Impl, broker: broker}, nil
}

func (p *CommandPluginImpl) Client(broker *plugin.MuxBroker, c *plugin.Client) (interface{}, error) {
	return &CommandPluginRPCClient{client: c, broker: broker}, nil
}

// CommandPluginRPCServer is the RPC server that CommandPluginRPC talks to
type CommandPluginRPCServer struct {
	Impl   CommandPlugin
	broker *plugin.MuxBroker
}

// Name implements the server side of the RPC interface
func (s *CommandPluginRPCServer) Name(args interface{}, resp *string) error {
	*resp = s.Impl.Name()
	return nil
}

// Version implements the server side of the RPC interface
func (s *CommandPluginRPCServer) Version(args interface{}, resp *string) error {
	*resp = s.Impl.Version()
	return nil
}

// Execute implements the server side of the RPC interface
func (s *CommandPluginRPCServer) Execute(args map[string]interface{}, resp *string) error {
	result, err := s.Impl.Execute(args)
	*resp = result
	return err
}

// GetCapabilities implements the server side of the RPC interface
func (s *CommandPluginRPCServer) GetCapabilities(args interface{}, resp *[]string) error {
	*resp = s.Impl.GetCapabilities()
	return nil
}

// CommandPluginRPCClient is the client implementation
type CommandPluginRPCClient struct {
	client *plugin.Client
	broker *plugin.MuxBroker
}

// Name calls the plugin's Name method via RPC
func (c *CommandPluginRPCClient) Name() string {
	var resp string
	err := c.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		return ""
	}
	return resp
}

// Version calls the plugin's Version method via RPC
func (c *CommandPluginRPCClient) Version() string {
	var resp string
	err := c.client.Call("Plugin.Version", new(interface{}), &resp)
	if err != nil {
		return ""
	}
	return resp
}

// Execute calls the plugin's Execute method via RPC
func (c *CommandPluginRPCClient) Execute(args map[string]interface{}) (string, error) {
	var resp string
	err := c.client.Call("Plugin.Execute", args, &resp)
	return resp, err
}

// GetCapabilities calls the plugin's GetCapabilities method via RPC
func (c *CommandPluginRPCClient) GetCapabilities() []string {
	var resp []string
	err := c.client.Call("Plugin.GetCapabilities", new(interface{}), &resp)
	if err != nil {
		return []string{}
	}
	return resp
}