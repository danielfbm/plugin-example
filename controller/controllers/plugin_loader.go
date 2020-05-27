package controllers

import (
	ext "github.com/danielfbm/plugin-example/controller/extension"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-plugin"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginLoader struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	PluginFolder string
	Manager      ext.Manager
}

func (r *PluginLoader) Start(stopChan <-chan struct{}) error {
	pluginsList, err := plugin.Discover("*.po", r.PluginFolder)
	if err != nil {
		return err
	}
	if len(pluginsList) > 0 {
		for _, p := range pluginsList {
			if p == "" {

			}
		}
	}
	<-stopChan
	return nil
}
