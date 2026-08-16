package main

import (
	"syscall"

	"github.com/rolldone/donisgo/internal/console"
	"github.com/rolldone/donisgo/internal/plugins"
	donisfinance "github.com/rolldone/donisgo/plugins/donisfinance"
)

func main() {
	// Ensure console commands create group-writable files by default
	syscall.Umask(0o002)

	// Register donisfinance plugin so its console commands are available
	console.RegisterAdditionalPlugins([]plugins.Plugin{donisfinance.New()})
	console.Execute()
}
