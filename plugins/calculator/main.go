package main

import (
	ldagent "github.com/your-org/ld-agent-go"
)

// AddNumbers adds two numbers together
func AddNumbers(a, b float64) float64 {
	return a + b
}

// ModuleInfo contains the plugin metadata
var ModuleInfo = ldagent.ModuleInfo{
	Name:        "Simple Calculator",
	Description: "Basic arithmetic operations",
	Author:      "ld-agent Team",
	Version:     "1.0.0",
	Platform:    "any",
	GoRequires:  ">=1.21",
	Dependencies: []string{},
	EnvironmentVariables: map[string]ldagent.EnvVar{},
}

// ModuleExports defines what this plugin exports
var ModuleExports = ldagent.ModuleExports{
	Tools: []ldagent.Tool{
		{
			Name:        "add_numbers",
			Description: "Add two numbers together and return the result",
			Function:    AddNumbers,
			Parameters: map[string]ldagent.Parameter{
				"a": {
					Type:        "float64",
					Description: "First number to add",
					Required:    true,
				},
				"b": {
					Type:        "float64", 
					Description: "Second number to add",
					Required:    true,
				},
			},
			ReturnType: "float64",
		},
	},
}

// Init is called when the plugin is loaded (optional)
func Init() error {
	// Any initialization logic here
	return nil
} 
