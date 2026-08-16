package pluginloader

import (
	"github.com/rolldone/donisgo/internal/plugins"
)

// RegisterCorePlugins registers the set of core plugins supported by donisgo.
func RegisterCorePlugins() {
	plugins.RegisterPlugins([]plugins.Plugin{
		// Add core plugins here
		// e.g., coreplugin.New(),
	})
}
