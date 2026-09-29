package main

import (
	"permissions-storage/internal"

	"github.com/untrustedmodders/go-plugify"
)

const pluginName = "permissions_storage"

func init() {
	plugin := internal.NewPermissionsStoragePlugin()
	plugin.Plugin = plugify.NewPlugin(pluginName, plugin.OnPluginStart, plugin.OnPluginUpdate, plugin.OnPluginEnd)
}

func main() {}
