package main

import (
	"syscall"

	"go_framework/internal/console"
	"go_framework/internal/plugins"
	donisfinance "go_framework/plugins/donisfinance"
)

func main() {
	// Ensure console commands create group-writable files by default
	syscall.Umask(0o002)

	// Register donisfinance plugin so its console commands are available
	console.RegisterAdditionalPlugins([]plugins.Plugin{donisfinance.New()})
	console.Execute()
}
