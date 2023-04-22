package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"plugin"
	"time"

	"github.com/mamur-rezeki/gowhatsplugins/types"
)

var PluginRoot = "plugins"
var PluginRootMod time.Time = time.Time{}
var GoWhatsPlugins = map[string]*types.Plugin{}

func loadPlugins() {

	GoWhatsPlugins = map[string]*types.Plugin{}

	if sd, err := os.Stat(PluginRoot); err != nil {
		log.Println("Creating plugin dir", PluginRoot)
		os.Mkdir(PluginRoot, fs.ModeAppend)
	} else {
		PluginRootMod = sd.ModTime()
	}

	if entries, err := os.ReadDir(PluginRoot); err == nil {
		for _, entry := range entries {
			if path.Ext(entry.Name()) == ".so" {
				var plugin_path = path.Join(PluginRoot, entry.Name())
				if tmpp, err := plugin.Open(plugin_path); err == nil {
					if plug, err := tmpp.Lookup("Plugin"); err == nil {

						GoWhatsPlugins[entry.Name()] = plug.(*types.Plugin)
					}
				} else {
					log.Println(err)
				}
			}
		}
	}
}
