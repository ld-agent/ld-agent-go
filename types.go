package ldagent

import (
	"reflect"
)

// EnvVar represents an environment variable configuration
type EnvVar struct {
	Description string `json:"description" validate:"required"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

// ModuleInfo contains plugin metadata
type ModuleInfo struct {
	Name                 string             `json:"name" validate:"required"`
	Description          string             `json:"description" validate:"required"`
	Author               string             `json:"author" validate:"required"`
	Version              string             `json:"version" validate:"required"`
	Platform             string             `json:"platform" validate:"required"`
	PythonRequires       string             `json:"python_requires,omitempty"`
	GoRequires           string             `json:"go_requires,omitempty"`
	Dependencies         []string           `json:"dependencies"`
	EnvironmentVariables map[string]EnvVar  `json:"environment_variables"`
}

// Tool represents a callable function/tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Function    interface{}            `json:"-"` // The actual Go function
	Parameters  map[string]Parameter   `json:"parameters"`
	ReturnType  string                 `json:"return_type"`
}

// Parameter represents a function parameter
type Parameter struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
}

// ModuleExports defines what a plugin exports
type ModuleExports struct {
	Tools     []Tool        `json:"tools"`
	Agents    []interface{} `json:"agents,omitempty"`
	Resources []interface{} `json:"resources,omitempty"`
	Models    []interface{} `json:"models,omitempty"`
	Utilities []Tool        `json:"utilities,omitempty"`
}

// Plugin represents a loaded plugin
type Plugin struct {
	Name     string        `json:"name"`
	Info     ModuleInfo    `json:"info"`
	Exports  ModuleExports `json:"exports"`
	FilePath string        `json:"file_path"`
}

// ToolRegistry manages all loaded tools
type ToolRegistry struct {
	Tools    map[string]Tool    `json:"tools"`
	Plugins  map[string]Plugin  `json:"plugins"`
	Metadata map[string]ModuleInfo `json:"metadata"`
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		Tools:    make(map[string]Tool),
		Plugins:  make(map[string]Plugin),
		Metadata: make(map[string]ModuleInfo),
	}
}

// RegisterTool registers a tool in the registry
func (tr *ToolRegistry) RegisterTool(pluginName string, tool Tool) {
	toolName := pluginName + "." + tool.Name
	tr.Tools[toolName] = tool
}

// GetTool retrieves a tool by name
func (tr *ToolRegistry) GetTool(name string) (Tool, bool) {
	tool, exists := tr.Tools[name]
	return tool, exists
}

// ListTools returns all tool names
func (tr *ToolRegistry) ListTools() []string {
	var tools []string
	for name := range tr.Tools {
		tools = append(tools, name)
	}
	return tools
}

// CallTool calls a tool with the given arguments
func (tr *ToolRegistry) CallTool(name string, args map[string]interface{}) (interface{}, error) {
	tool, exists := tr.Tools[name]
	if !exists {
		return nil, ErrToolNotFound
	}

	// Use reflection to call the function
	fn := reflect.ValueOf(tool.Function)
	if fn.Kind() != reflect.Func {
		return nil, ErrInvalidTool
	}

	// For now, we'll implement a simple calling mechanism
	// In a full implementation, you'd want proper parameter mapping
	fnType := fn.Type()
	var callArgs []reflect.Value

	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		// Simple parameter mapping - in practice you'd want more sophisticated handling
		if i < len(args) {
			// Convert args to proper types - simplified for demo
			arg := reflect.ValueOf(getArgByIndex(args, i))
			if arg.Type().ConvertibleTo(paramType) {
				callArgs = append(callArgs, arg.Convert(paramType))
			} else {
				callArgs = append(callArgs, reflect.Zero(paramType))
			}
		} else {
			callArgs = append(callArgs, reflect.Zero(paramType))
		}
	}

	results := fn.Call(callArgs)
	if len(results) > 0 {
		return results[0].Interface(), nil
	}
	return nil, nil
}

// Helper function to get argument by index (simplified)
func getArgByIndex(args map[string]interface{}, index int) interface{} {
	i := 0
	for _, v := range args {
		if i == index {
			return v
		}
		i++
	}
	return nil
} 
