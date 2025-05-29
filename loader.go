package ldagent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

// Loader manages plugin discovery and loading
type Loader struct {
	PluginsDir string
	Silent     bool
	Registry   *ToolRegistry
}

// NewLoader creates a new plugin loader
func NewLoader(pluginsDir string, silent bool) *Loader {
	return &Loader{
		PluginsDir: pluginsDir,
		Silent:     silent,
		Registry:   NewToolRegistry(),
	}
}

// log prints a message if not in silent mode
func (l *Loader) log(message string) {
	if !l.Silent {
		log.Printf("🔌 %s", message)
	}
}

// isCompatible checks if a plugin is compatible with the current environment
func (l *Loader) isCompatible(info ModuleInfo) bool {
	// Check platform
	if info.Platform != "any" && info.Platform != runtime.GOOS {
		return false
	}
	
	// Check Go version (basic check)
	if info.GoRequires != "" {
		// In a full implementation, you'd parse and compare versions
		// For now, we'll just check if it's specified
		if !strings.HasPrefix(info.GoRequires, ">=") {
			return false
		}
	}
	
	return true
}

// LoadPlugin loads a single plugin file (.so)
func (l *Loader) LoadPlugin(pluginPath string) error {
	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		l.log(fmt.Sprintf("❌ Failed to load %s: %v", filepath.Base(pluginPath), err))
		return fmt.Errorf("%w: %v", ErrPluginLoadFailed, err)
	}
	
	// Get module info
	infoSym, err := p.Lookup("ModuleInfo")
	if err != nil {
		l.log(fmt.Sprintf("⚠️  %s missing ModuleInfo", filepath.Base(pluginPath)))
		return ErrMissingMetadata
	}
	
	moduleInfo, ok := infoSym.(*ModuleInfo)
	if !ok {
		l.log(fmt.Sprintf("⚠️  %s invalid ModuleInfo type", filepath.Base(pluginPath)))
		return ErrMissingMetadata
	}
	
	// Get module exports
	exportsSym, err := p.Lookup("ModuleExports")
	if err != nil {
		l.log(fmt.Sprintf("⚠️  %s missing ModuleExports", filepath.Base(pluginPath)))
		return ErrMissingExports
	}
	
	moduleExports, ok := exportsSym.(*ModuleExports)
	if !ok {
		l.log(fmt.Sprintf("⚠️  %s invalid ModuleExports type", filepath.Base(pluginPath)))
		return ErrMissingExports
	}
	
	// Check compatibility
	if !l.isCompatible(*moduleInfo) {
		l.log(fmt.Sprintf("⚠️  %s not compatible", filepath.Base(pluginPath)))
		return ErrIncompatiblePlugin
	}
	
	// Register the plugin
	pluginName := strings.TrimSuffix(filepath.Base(pluginPath), ".so")
	
	plugin := Plugin{
		Name:     pluginName,
		Info:     *moduleInfo,
		Exports:  *moduleExports,
		FilePath: pluginPath,
	}
	
	l.Registry.Plugins[pluginName] = plugin
	l.Registry.Metadata[pluginName] = *moduleInfo
	
	// Register tools
	for _, tool := range moduleExports.Tools {
		l.Registry.RegisterTool(pluginName, tool)
	}
	
	// Call initialization function if provided
	if initSym, err := p.Lookup("Init"); err == nil {
		if initFunc, ok := initSym.(func() error); ok {
			if err := initFunc(); err != nil {
				l.log(fmt.Sprintf("⚠️  Failed to initialize %s: %v", pluginName, err))
			} else {
				l.log(fmt.Sprintf("🔧 Initialized %s", pluginName))
			}
		}
	}
	
	l.log(fmt.Sprintf("✅ Loaded %s %s", moduleInfo.Name, moduleInfo.Version))
	return nil
}

// LoadAll loads all plugins from the plugins directory
func (l *Loader) LoadAll() int {
	if _, err := os.Stat(l.PluginsDir); os.IsNotExist(err) {
		return 0
	}
	
	loaded := 0
	
	// Load .so files (compiled Go plugins)
	err := filepath.Walk(l.PluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(path, ".so") {
			if err := l.LoadPlugin(path); err == nil {
				loaded++
			}
		}
		
		return nil
	})
	
	if err != nil {
		l.log(fmt.Sprintf("Error scanning plugins directory: %v", err))
	}
	
	if !l.Silent && loaded > 0 {
		l.log(fmt.Sprintf("Loaded %d plugins", loaded))
	}
	
	return loaded
}

// GetTool retrieves a tool by name
func (l *Loader) GetTool(name string) (Tool, bool) {
	return l.Registry.GetTool(name)
}

// ListTools returns all available tool names
func (l *Loader) ListTools() []string {
	return l.Registry.ListTools()
}

// ListPlugins returns all loaded plugins with metadata
func (l *Loader) ListPlugins() map[string]ModuleInfo {
	return l.Registry.Metadata
}

// CallTool calls a tool with the given arguments
func (l *Loader) CallTool(name string, args map[string]interface{}) (interface{}, error) {
	return l.Registry.CallTool(name, args)
} 
