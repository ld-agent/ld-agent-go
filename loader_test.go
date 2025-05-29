package ldagent

import (
	"testing"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader("plugins", true)
	
	if loader == nil {
		t.Fatal("NewLoader returned nil")
	}
	
	if loader.PluginsDir != "plugins" {
		t.Errorf("Expected plugins dir 'plugins', got '%s'", loader.PluginsDir)
	}
	
	if !loader.Silent {
		t.Error("Expected silent mode to be true")
	}
	
	if loader.Registry == nil {
		t.Error("Registry should not be nil")
	}
}

func TestToolRegistry(t *testing.T) {
	registry := NewToolRegistry()
	
	if registry == nil {
		t.Fatal("NewToolRegistry returned nil")
	}
	
	// Test empty registry
	tools := registry.ListTools()
	if len(tools) != 0 {
		t.Errorf("Expected 0 tools, got %d", len(tools))
	}
	
	// Test tool registration
	testTool := Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Function:    func(a, b int) int { return a + b },
		Parameters: map[string]Parameter{
			"a": {Type: "int", Description: "First number", Required: true},
			"b": {Type: "int", Description: "Second number", Required: true},
		},
		ReturnType: "int",
	}
	
	registry.RegisterTool("test_plugin", testTool)
	
	// Check if tool was registered
	tools = registry.ListTools()
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
	
	if tools[0] != "test_plugin.test_tool" {
		t.Errorf("Expected tool name 'test_plugin.test_tool', got '%s'", tools[0])
	}
	
	// Test tool retrieval
	tool, exists := registry.GetTool("test_plugin.test_tool")
	if !exists {
		t.Error("Tool should exist")
	}
	
	if tool.Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", tool.Name)
	}
}

func TestModuleInfoCompatibility(t *testing.T) {
	loader := NewLoader("plugins", true)
	
	// Test compatible module
	compatibleInfo := ModuleInfo{
		Name:        "Test Plugin",
		Description: "A test plugin",
		Author:      "Test Author",
		Version:     "1.0.0",
		Platform:    "any",
		GoRequires:  ">=1.21",
	}
	
	if !loader.isCompatible(compatibleInfo) {
		t.Error("Module should be compatible")
	}
	
	// Test incompatible platform
	incompatibleInfo := ModuleInfo{
		Name:        "Test Plugin",
		Description: "A test plugin", 
		Author:      "Test Author",
		Version:     "1.0.0",
		Platform:    "nonexistent_platform",
		GoRequires:  ">=1.21",
	}
	
	if loader.isCompatible(incompatibleInfo) {
		t.Error("Module should not be compatible")
	}
} 
