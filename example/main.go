package main

import (
	"fmt"
	"log"
	
	ldagent "github.com/your-org/ld-agent-go"
)

func main() {
	fmt.Println("🔗 ld-agent Go Example")
	fmt.Println("Dynamic linking for agentic systems")
	fmt.Println()

	// Create a new loader
	loader := ldagent.NewLoader("../plugins", false)
	
	// Load all plugins
	loaded := loader.LoadAll()
	fmt.Printf("📦 Loaded %d plugins\n\n", loaded)
	
	// List available tools
	tools := loader.ListTools()
	fmt.Printf("🔧 Available tools (%d):\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("   • %s\n", tool)
	}
	fmt.Println()
	
	// List loaded plugins
	plugins := loader.ListPlugins()
	fmt.Printf("🔌 Loaded plugins (%d):\n", len(plugins))
	for _, info := range plugins {
		fmt.Printf("   • %s v%s - %s\n", info.Name, info.Version, info.Description)
	}
	fmt.Println()
	
	// Test the calculator tool
	if len(tools) > 0 {
		fmt.Println("🧮 Testing calculator.add_numbers:")
		
		// Call the add_numbers tool
		args := map[string]interface{}{
			"a": 15.5,
			"b": 24.3,
		}
		
		result, err := loader.CallTool("calculator.add_numbers", args)
		if err != nil {
			log.Printf("❌ Error calling tool: %v", err)
		} else {
			fmt.Printf("   15.5 + 24.3 = %v\n", result)
		}
		
		// Test with different numbers
		args2 := map[string]interface{}{
			"a": 100.0,
			"b": 200.0,
		}
		
		result2, err := loader.CallTool("calculator.add_numbers", args2)
		if err != nil {
			log.Printf("❌ Error calling tool: %v", err)
		} else {
			fmt.Printf("   100.0 + 200.0 = %v\n", result2)
		}
	}
	
	fmt.Println()
	fmt.Println("✅ ld-agent Go example completed!")
} 
