# ld-agent Go

**Dynamic linking for agentic systems in Go**

ld-agent-go is the Go implementation of ld-agent, a dynamic linker for AI capabilities. Just like `ld-so` discovers and links shared libraries at runtime, ld-agent discovers and links agentic capabilities to create composable AI systems.

## Quick Start

### 1. Build and Run the Example

```bash
# Build the calculator plugin and example
make run
```

This will:
- Build the calculator plugin as a shared library (`calculator.so`)
- Build the example program
- Run the example, which loads the plugin and tests the `add_numbers` tool

### 2. Expected Output

```
🔗 ld-agent Go Example
Dynamic linking for agentic systems

🔌 ✅ Loaded Simple Calculator 1.0.0
📦 Loaded 1 plugins

🔧 Available tools (1):
   • calculator.add_numbers

🔌 Loaded plugins (1):
   • Simple Calculator v1.0.0 - Basic arithmetic operations

🧮 Testing calculator.add_numbers:
   15.5 + 24.3 = 39.8
   100.0 + 200.0 = 300

✅ ld-agent Go example completed!
```

## Creating Plugins

### Plugin Structure

Go plugins must be built with `-buildmode=plugin` and export specific symbols:

```go
package main

import ldagent "github.com/your-org/ld-agent-go"

// Your tool function
func AddNumbers(a, b float64) float64 {
    return a + b
}

// Required: Plugin metadata
var ModuleInfo = &ldagent.ModuleInfo{
    Name:        "Simple Calculator",
    Description: "Basic arithmetic operations", 
    Author:      "Your Name",
    Version:     "1.0.0",
    Platform:    "any",
    GoRequires:  ">=1.21",
    Dependencies: []string{},
    EnvironmentVariables: map[string]ldagent.EnvVar{},
}

// Required: Plugin exports
var ModuleExports = &ldagent.ModuleExports{
    Tools: []ldagent.Tool{
        {
            Name:        "add_numbers",
            Description: "Add two numbers together",
            Function:    AddNumbers,
            Parameters: map[string]ldagent.Parameter{
                "a": {Type: "float64", Description: "First number", Required: true},
                "b": {Type: "float64", Description: "Second number", Required: true},
            },
            ReturnType: "float64",
        },
    },
}

// Optional: Initialization function
func Init() error {
    // Plugin initialization logic
    return nil
}
```

### Building Plugins

```bash
# Build as a shared library
go build -buildmode=plugin -o plugins/myplugin.so ./path/to/plugin/

# Or use the Makefile
make build-plugin
```

## Using ld-agent in Your Code

```go
package main

import (
    "fmt"
    ldagent "github.com/your-org/ld-agent-go"
)

func main() {
    // Create loader
    loader := ldagent.NewLoader("plugins", false)
    
    // Load all plugins
    loaded := loader.LoadAll()
    fmt.Printf("Loaded %d plugins\n", loaded)
    
    // List available tools
    tools := loader.ListTools()
    for _, tool := range tools {
        fmt.Printf("Tool: %s\n", tool)
    }
    
    // Call a tool
    args := map[string]interface{}{
        "a": 10.0,
        "b": 20.0,
    }
    
    result, err := loader.CallTool("calculator.add_numbers", args)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %v\n", result)
    }
}
```

## API Reference

### Core Types

- `ModuleInfo` - Plugin metadata (name, version, dependencies, etc.)
- `ModuleExports` - What the plugin exports (tools, agents, resources)
- `Tool` - A callable function with metadata
- `Parameter` - Function parameter definition
- `Loader` - Main plugin loader and registry

### Key Functions

- `NewLoader(pluginsDir, silent)` - Create a new plugin loader
- `LoadAll()` - Load all plugins from directory
- `LoadPlugin(path)` - Load a specific plugin
- `GetTool(name)` - Get a tool by name
- `CallTool(name, args)` - Call a tool with arguments
- `ListTools()` - List all available tools
- `ListPlugins()` - List all loaded plugins

## Build Commands

```bash
# Build everything
make all

# Build just the plugin
make build-plugin

# Build just the example
make build-example

# Run the example
make run

# Run tests
make test

# Clean build artifacts
make clean

# Initialize Go modules
make init
```

## Requirements

- Go 1.21 or later
- Unix-like system (Linux, macOS) - Go plugins require CGO and don't work on Windows

## Architecture

ld-agent-go follows the same conceptual model as the Python version:

1. **Discovery** - Scans `plugins/` directory for `.so` files
2. **Loading** - Uses Go's `plugin` package to load shared libraries
3. **Symbol Resolution** - Looks up `ModuleInfo` and `ModuleExports` symbols
4. **Registration** - Registers tools in a global registry
5. **Execution** - Uses reflection to call tools with proper arguments

This enables truly modular AI systems where capabilities can be mixed and matched at runtime without recompilation.

## Limitations

- Go plugins only work on Unix-like systems (Linux, macOS)
- Plugins must be compiled with the same Go version as the loader
- Plugin functions must use basic Go types for parameters (no complex structs across plugin boundaries)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Run `make test` to ensure tests pass
5. Submit a pull request

## License

Same as the main ld-agent project. 
