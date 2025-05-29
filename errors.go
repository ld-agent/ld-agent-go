package ldagent

import "errors"

var (
	// ErrToolNotFound is returned when a requested tool is not found
	ErrToolNotFound = errors.New("tool not found")
	
	// ErrInvalidTool is returned when a tool is not a valid function
	ErrInvalidTool = errors.New("invalid tool: not a function")
	
	// ErrPluginLoadFailed is returned when a plugin fails to load
	ErrPluginLoadFailed = errors.New("plugin load failed")
	
	// ErrIncompatiblePlugin is returned when a plugin is not compatible
	ErrIncompatiblePlugin = errors.New("plugin not compatible with current environment")
	
	// ErrMissingMetadata is returned when plugin metadata is missing
	ErrMissingMetadata = errors.New("plugin missing required metadata")
	
	// ErrMissingExports is returned when plugin exports are missing
	ErrMissingExports = errors.New("plugin missing exports")
) 
