package plugin

type PluginsData = map[string]interface{};

type Plugin interface {
	LoadPlugin(pluginsData *PluginsData, devMode *bool)
}