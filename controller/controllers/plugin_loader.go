package controllers

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	pluginsv1alpha1 "github.com/danielfbm/plugin-example/controller/api/v1alpha1"
	ext "github.com/danielfbm/plugin-example/controller/extension"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var LocalPluginFinilizer = "local-plugin-finilizer"

type PluginLoader struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	PluginFolder string
	Manager      ext.Manager
}

func (r *PluginLoader) Start(stopChan <-chan struct{}) error {
	r.Log.Info("Starting...", "folder", r.PluginFolder)
	pluginsList, err := plugin.Discover("*.po", r.PluginFolder)
	if err != nil {
		return err
	}
	r.Log.Info("Found local plugins", "plugins", pluginsList)
	if len(pluginsList) > 0 {

		opts := ext.PluginLoadOptions{
			Config: &plugin.ClientConfig{
				Logger: hclog.Default(),
			},
		}
		for _, p := range pluginsList {
			_, filename := filepath.Split(p)
			filename = strings.TrimSuffix(filename, ".po")
			pOpts := opts.Copy()
			pOpts.LocalPluginPath = p
			r.Log.V(0).Info("Loading plugin", "name", filename, "path", p)
			err := r.Manager.Load(filename, pOpts)
			pluginInstance := r.newPlugin(filename, p, err)
			if err = r.upsertPlugin(context.Background(), pluginInstance); err != nil {
				r.Log.Error(err, "Plugin upsert error", "name", pluginInstance.Name)
			}
			// r.Client.Update(pluginInstance)
			// create CRD

		}
	}
	r.Log.Info("Load finished...")
	<-stopChan
	return nil
}

func (r *PluginLoader) newPlugin(name, path string, err error) (res *pluginsv1alpha1.Plugin) {
	load := err == nil
	message, reason := "", ""
	if err != nil {
		message = err.Error()
		reason = "LoadError"
	}
	res = &pluginsv1alpha1.Plugin{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  "default",
			Finalizers: []string{LocalPluginFinilizer},
		},
		Spec: pluginsv1alpha1.PluginSpec{
			Local: &pluginsv1alpha1.LocalPluginSpec{
				Path: path,
			},
		},
		Status: pluginsv1alpha1.PluginStatus{
			Conditions: []pluginsv1alpha1.Condition{
				pluginsv1alpha1.Condition{
					Type:    "Load",
					Ready:   strconv.FormatBool(load),
					Message: message,
					Reason:  reason,
				},
			},
		},
	}
	return
}

func (r *PluginLoader) upsertPlugin(ctx context.Context, plug *pluginsv1alpha1.Plugin) (err error) {
	existing := &pluginsv1alpha1.Plugin{}
	if foundErr := r.Client.Get(ctx, client.ObjectKey{Name: plug.Name, Namespace: plug.Namespace}, existing); foundErr != nil && errors.IsNotFound(foundErr) {
		// create
		err = r.Client.Create(ctx, plug)
	} else {
		// update
		existing.Spec = plug.Spec
		existing.Status = plug.Status
		err = r.Client.Update(ctx, existing)
	}
	return
}
